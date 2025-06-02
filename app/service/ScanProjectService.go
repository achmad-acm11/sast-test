package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"sast-integration/app/dbo/entity"
	"sast-integration/app/dto"
	"sast-integration/app/dto/request"
	"sast-integration/app/helper"
	"strings"
)

type projectPathInfo struct {
	scannedProjectFilePath string
	projectPath            string
	typeVersion            string
}

func (p ProjectServiceImpl) scanningRepository(ctx *gin.Context, project entity.Project) {
	p.stdLog.InfoFunction(fmt.Sprintf("START SCANNING REPOSITORY %s", project.Key))

	pathData := p.preparingProjectPath(ctx, project)
	p.prepareRunScanAndAction(ctx, project, pathData)
	p.saveResultToDb(ctx, project, pathData)

	p.stdLog.InfoFunction(fmt.Sprintf("END SCANNING REPOSITORY %s", project.Key))
}

// PREPARING PROJECT PATH PHASE
func (p ProjectServiceImpl) preparingProjectPath(ctx *gin.Context, project entity.Project) projectPathInfo {
	var curDir, _ = os.Getwd()
	typeVersion := "repository"

	// Example Result of Project Path => /_project-repository/workspace/test_project-sca
	RepositoryPath := fmt.Sprintf("%s%s", curDir, p.prefixProjectFolder)
	projectPath := fmt.Sprintf("%s%s-sast-semgrep", RepositoryPath, project.Key)

	// Example Result of Scanned Project File Path => /_scanned-project-files/scan_repo_file_test_project_repository-.json
	ProjectFilePath := fmt.Sprintf("%s%s", curDir, p.prefixResultFolder)
	projectFileName := "scan_repo_file_" + project.Key + "_" + typeVersion + "-sast-semgrep" + ".json"
	scannedProjectFilePath := ProjectFilePath + projectFileName

	data := projectPathInfo{
		scannedProjectFilePath: scannedProjectFilePath,
		projectPath:            projectPath,
		typeVersion:            typeVersion,
	}

	return data
}

// PREPARING RUN SCAN AND ACTION PHASE
func (p ProjectServiceImpl) prepareRunScanAndAction(ctx *gin.Context, project entity.Project, pathData projectPathInfo) {

	out, err := p.runScanCommand(pathData)

	if err != nil {
		p.checkErrorScan(ctx, project, err)
	}

	p.stdLog.InfoFunction("Success scan command: " + fmt.Sprintf("%q\n", out))
}

func (p ProjectServiceImpl) runScanCommand(pathData projectPathInfo) ([]byte, error) {
	semgrepCli := p.semgrepCli.Init().AddProjectPath(pathData.projectPath).AddOutput(pathData.scannedProjectFilePath)
	out, err := semgrepCli.Exec()

	return out, err
}

func (p ProjectServiceImpl) checkErrorScan(ctx *gin.Context, project entity.Project, err error) {
	if err != nil {
		p.failedScan(ctx, project)
		helper.ErrorHandler(err)
	}
}

// SAVE RESULT TO DB PHASE
func (p ProjectServiceImpl) saveResultToDb(ctx *gin.Context, project entity.Project, pathData projectPathInfo) {
	scanVersion := p.incrementScanVersion(ctx, project)

	results := p.mappingResultOutput(pathData.scannedProjectFilePath, project.Id, project.Key)

	nothingToChange := p.compareNewResultWithPrevResult(ctx, project.Id, results, scanVersion)

	if nothingToChange {
		project.StatusScan = 3
		//project.Message = "Nothing To Update"
		p.repo.Update(ctx, p.db, project)

		p.stdLog.InfoFunction("Done scanning project repository... (nothing to update)")
		return
	}

	for _, result := range results {
		resultRule := p.repoResult.GetOneResultByProjectIdAndRule(ctx, p.db, project.Id, result.Rule)
		p.addOrUpdateResult(ctx, result, resultRule, scanVersion, pathData.typeVersion)
	}

	p.deleteAllPrevVersionResult(ctx, project, scanVersion)

	project.StatusScan = 3
	//project.StatusMessage = ""
	project.CurrentScanVersion = scanVersion
	p.repo.Update(ctx, p.db, project)

	issues := []request.IssueRequest{}
	for _, result := range results {
		issues = append(issues, request.IssueRequest{
			Title:        result.Title,
			Rule:         result.Rule,
			Path:         result.Path,
			Line:         result.Line,
			Type:         result.Type,
			Description:  result.Description,
			Severity:     result.Severity,
			References:   result.References,
			LastFoundAt:  result.LastFoundAt,
			StatusResult: result.StatusResult,
		})
	}

	p.asocApi.SendResult(request.AsocSendResultRequest{
		ProjectKey:  project.Key,
		ScanVersion: scanVersion,
		CreatedAt:   "",
		Issues:      issues,
	})
}

func (p ProjectServiceImpl) incrementScanVersion(ctx *gin.Context, project entity.Project) int {
	lastScanResult := p.repoResult.GetLastByProjectId(ctx, p.db, project.Id)

	scanVersion := 0

	if lastScanResult.Id > 0 {
		if lastScanResult.ScanVersion <= project.CurrentScanVersion {
			scanVersion = project.CurrentScanVersion + 1
		} else {
			scanVersion = int(lastScanResult.ScanVersion) + 1
		}
	} else {
		scanVersion++
	}

	return scanVersion
}

func (p ProjectServiceImpl) mappingResultOutput(filePathResult string, projectId int, projectKey string) []entity.Result {
	var dataJson dto.SemgrepResult

	sourceFile, _ := os.Open(filePathResult)
	defer sourceFile.Close()

	byteValue1, _ := ioutil.ReadAll(sourceFile)
	_ = json.Unmarshal(byteValue1, &dataJson)

	results := []entity.Result{}

	for _, data := range dataJson.Results {
		var rules []string
		var titles []string
		extra := data.Extra
		metadata := extra.Metadata
		for _, cwe := range data.Extra.Metadata.CWE {
			split := strings.Split(cwe, ":")
			rules = append(rules, strings.TrimSpace(split[0]))
			titles = append(titles, strings.TrimSpace(split[1]))
		}

		tmpResult := new(entity.Result)
		tmpResult.ProjectId = projectId
		tmpResult.ProjectKey = projectKey
		tmpResult.Rule = strings.Join(rules, "|")
		tmpResult.References = strings.Join(metadata.References, "|")
		tmpResult.Title = strings.Join(titles, "|")
		tmpResult.Description = extra.Message
		tmpResult.Severity = metadata.Impact
		tmpResult.Type = metadata.Subcategory[0]

		results = append(results, *tmpResult)
	}
	return results
}

func (p ProjectServiceImpl) compareNewResultWithPrevResult(ctx *gin.Context, projectId int, newResults []entity.Result, currentVersion int) bool {
	var fixedVulnerabilityBool = false
	var newVulnerabilityBool = false

	prevResults := p.repoResult.GetAllByProjectIdAndScanVersion(ctx, p.db, projectId, currentVersion-1)
	newVulnerabilityBool = p.checkingNotHaveSameRule(prevResults, newResults)
	fixedVulnerabilityBool = p.checkingNotHaveSameRule(prevResults, prevResults)

	if fixedVulnerabilityBool == false && newVulnerabilityBool == false {
		return true
	}
	return false
}

func (p ProjectServiceImpl) checkingNotHaveSameRule(result1 []entity.Result, result2 []entity.Result) bool {
	if len(result1) == 0 {
		return true
	}

	ruleSet := make(map[string]bool)

	for _, result := range result1 {
		if _, ok := ruleSet[result.Rule]; !ok {
			ruleSet[result.Rule] = true
		}
	}

	for _, result := range result2 {
		if _, ok := ruleSet[result.Rule]; ok {
			continue
		}
		return true
	}
	return false
}

func (p ProjectServiceImpl) addOrUpdateResult(ctx *gin.Context, result entity.Result, oldResult entity.Result, scanVersion int, typeVersion string) {
	if oldResult.Id != 0 {
		tmpScanVersion := oldResult.ScanVersion
		tmpLastUpdate := oldResult.UpdatedAt.Format("2006-01-02 15:04:05")

		oldResult.ScanVersion = scanVersion
		oldResult.StatusResult = 0
		oldResult.LastFoundAt = fmt.Sprintf("Scan no. %d at %s", tmpScanVersion, tmpLastUpdate)
		p.repoResult.Update(ctx, p.db, oldResult)
	} else {
		//result.ScanType = typeVersion
		result.ScanVersion = scanVersion
		result.LastFoundAt = ""
		p.repoResult.Create(ctx, p.db, result)
	}
}

func (p ProjectServiceImpl) deleteAllPrevVersionResult(ctx *gin.Context, project entity.Project, scanVersion int) {
	prevResults := p.repoResult.GetAllByProjectIdAndScanVersion(ctx, p.db, project.Id, scanVersion-1)
	for _, v := range prevResults {
		p.repoResult.DeleteOne(ctx, p.db, v)
	}
}

// UTILS SECTION
func (p ProjectServiceImpl) failedScan(ctx *gin.Context, project entity.Project) {
	project.StatusScan = 2
	// projectInformation.StatusMessage = statusMessage
	p.repo.Update(ctx, p.db, project)

	p.stdLog.WarningFunction(fmt.Sprintf("Project %v scan failed", project.Key))
}
