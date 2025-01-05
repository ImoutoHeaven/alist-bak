package tool

import (
	"context"
	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/task"
	"net/url"
	"path"
	"path/filepath"

	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DeletePolicy string

const (
	DeleteOnUploadSucceed DeletePolicy = "delete_on_upload_succeed"
	DeleteOnUploadFailed  DeletePolicy = "delete_on_upload_failed"
	DeleteNever           DeletePolicy = "delete_never"
	DeleteAlways          DeletePolicy = "delete_always"
)

type AddURLArgs struct {
	URL          string
	DstDirPath   string
	Tool         string
	DeletePolicy DeletePolicy
}

func AddURL(ctx context.Context, args *AddURLArgs) (task.TaskExtensionInfo, error) {
	// check storage
	storage, dstDirActualPath, err := op.GetStorageAndActualPath(args.DstDirPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed get storage")
	}
	// try putting url
	if args.Tool == "SimpleHttp" && tryPutUrl(ctx, storage, dstDirActualPath, args.URL) {
		return nil, nil
	}
	// check is it could upload
	if storage.Config().NoUpload {
		return nil, errors.WithStack(errs.UploadNotSupported)
	}
	// check path is valid
	obj, err := op.Get(ctx, storage, dstDirActualPath)
	if err != nil {
		if !errs.IsObjectNotFound(err) {
			return nil, errors.WithMessage(err, "failed get object")
		}
	} else {
		if !obj.IsDir() {
			// can't add to a file
			return nil, errors.WithStack(errs.NotFolder)
		}
	}

	// get tool
	tool, err := Tools.Get(args.Tool)
	if err != nil {
		return nil, errors.Wrapf(err, "failed get tool")
	}
	// check tool is ready
	if !tool.IsReady() {
		// try to init tool
		if _, err := tool.Init(); err != nil {
			return nil, errors.Wrapf(err, "failed init tool %s", args.Tool)
		}
	}

	uid := uuid.NewString()
	tempDir := filepath.Join(conf.Conf.TempDir, args.Tool, uid)
	deletePolicy := args.DeletePolicy

	switch args.Tool {
	case "115 Cloud":
		tempDir = args.DstDirPath
		// 防止将下载好的文件删除
		deletePolicy = DeleteNever
	case "pikpak":
		tempDir = args.DstDirPath
		// 防止将下载好的文件删除
		deletePolicy = DeleteNever
	case "thunder":
		tempDir = args.DstDirPath
		// 防止将下载好的文件删除
		deletePolicy = DeleteNever
	}

	taskCreator, _ := ctx.Value("user").(*model.User) // taskCreator is nil when convert failed
	t := &DownloadTask{
		TaskExtension: task.TaskExtension{
			Creator: taskCreator,
		},
		Url:          args.URL,
		DstDirPath:   args.DstDirPath,
		TempDir:      tempDir,
		DeletePolicy: deletePolicy,
		Toolname:     args.Tool,
		tool:         tool,
	}
	DownloadTaskManager.Add(t)
	return t, nil
}

func tryPutUrl(ctx context.Context, storage driver.Driver, dstDirActualPath, urlStr string) bool {
	_, ok := storage.(driver.PutURL)
	_, okResult := storage.(driver.PutURLResult)
	if !ok && !okResult {
		return false
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	dstName := path.Base(u.Path)
	err = op.PutURL(ctx, storage, dstDirActualPath, dstName, urlStr)
	return err == nil
}
