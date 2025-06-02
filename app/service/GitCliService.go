package service

import (
	"net/url"
	"sast-integration/app/dbo/cli"
	"sast-integration/app/dto/request"
	"sast-integration/app/helper"
	"sast-integration/app/shareVar"
)

type GitCliService interface {
	RunCloneRepo(request request.GitCloneRequest)
}

type GitCliServiceImpl struct {
	GitCli *cli.GitCli
	stdLog *helper.StandartLog
}

func NewGitCliService(gitCli *cli.GitCli) *GitCliServiceImpl {
	return &GitCliServiceImpl{
		GitCli: gitCli,
		stdLog: helper.NewStandardLog(shareVar.GitCli, shareVar.Service),
	}
}

func (g *GitCliServiceImpl) RunCloneRepo(request request.GitCloneRequest) {
	g.stdLog.NameFunc = "RunCloneRepo"
	g.stdLog.StartFunction(request)

	if request.Visibility == "PRIVATE" {
		request.RepoUrl = g.modifyUrlForPrivate(request.RepoUrl, request.Username, request.Token)
	}

	g.GitCli.InitClone(request.RepoUrl).
		AddParamBranch(request.BranchName).
		RepoDir(request.RepoName).
		Exec()

	g.stdLog.NameFunc = "RunCloneRepo"
	g.stdLog.EndFunction(nil)
}

func (g *GitCliServiceImpl) modifyUrlForPrivate(repoUrl string, username string, token string) string {
	parsedURL, err := url.Parse(repoUrl)
	if err != nil {
		panic(err)
	}

	parsedURL.User = url.UserPassword(username, token)

	return parsedURL.String()
}
