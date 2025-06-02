package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"sast-integration/app/controller"
	"sast-integration/app/dbo/cli"
	"sast-integration/app/dbo/repository"
	"sast-integration/app/dbo/repository/project"
	"sast-integration/app/service"
)

func ProjectRoute(router *gin.RouterGroup, db *gorm.DB, validate *validator.Validate) *gin.RouterGroup {
	projectRepository := repository.NewProjectRepository()
	projectAuthRepository := project.NewAuthRepository()
	resultRepository := repository.NewResultRepository()
	gitCli := cli.NewGitCli()
	semgrepCli := cli.NewSemgrepCli()

	gitCliService := service.NewGitCliService(gitCli)
	projectService := service.NewProjectService(projectRepository,
		projectAuthRepository,
		resultRepository,
		semgrepCli,
		gitCliService,
		validate,
		db)

	projectController := controller.NewProjectController(projectService)

	router.GET("projects", projectController.GetAllHandler)
	router.GET("project/:id", projectController.GetDetailByIdHandler)
	router.POST("project", projectController.CreateProjectHandler)
	router.POST("project/scan", projectController.ScanProjectHandler)
	router.PUT("project/:id", projectController.UpdateProjectHandler)
	router.DELETE("project/:id", projectController.DeleteProjectHandler)

	return router
}
