package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/seo/api/grpc/v1/gen"
	ctrl "github.com/JMURv/seo/internal/ctrl"
	"github.com/JMURv/seo/internal/dto"
	model "github.com/JMURv/seo/internal/models"
	utils "github.com/JMURv/seo/internal/models/mapper"
	"github.com/JMURv/seo/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHandler_GetSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()
	name, pk := "name", "pk"
	expectedSEO := &model.SEO{}

	t.Run(
		"Success", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			mockCtrl.EXPECT().GetSEO(gomock.Any(), name, pk).Return(expectedSEO, nil).Times(1)

			res, err := h.GetSEO(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"InvalidArgument", func(t *testing.T) {
			res, err := h.GetSEO(ctx, &pb.GetSEOReq{Name: "", Pk: ""})

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			mockCtrl.EXPECT().GetSEO(gomock.Any(), name, pk).Return(expectedSEO, ctrl.ErrNotFound).Times(1)

			res, err := h.GetSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			newErr := errors.New("new error")
			mockCtrl.EXPECT().GetSEO(gomock.Any(), name, pk).Return(expectedSEO, newErr).Times(1)

			res, err := h.GetSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)
}

func TestHandler_CreateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()

	req := &pb.SEOMsg{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "ogdescription",
		OGImage:       "ogimage",
		ObjName:       "objname",
		ObjPk:         "objpk",
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(
					&dto.CreateSEOResponse{
						Name: "name",
						PK:   "pk",
					}, nil,
				).
				Times(1)

			res, err := h.CreateSEO(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(nil, ctrl.ErrNotFound).
				Times(1)

			res, err := h.CreateSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				CreateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(nil, newErr).
				Times(1)

			res, err := h.CreateSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Title", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Description", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Keywords", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGTitle", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGDescription", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGImage", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OBJName", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OBJPK", func(t *testing.T) {
			res, err := h.CreateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)
}

func TestHandler_UpdateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()

	req := &pb.SEOMsg{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "ogdescription",
		OGImage:       "ogimage",
		ObjName:       "objname",
		ObjPk:         "objpk",
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(nil).
				Times(1)

			res, err := h.UpdateSEO(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(ctrl.ErrNotFound).
				Times(1)

			res, err := h.UpdateSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				UpdateSEO(gomock.Any(), utils.ProtoToModel(req)).
				Return(newErr).
				Times(1)

			res, err := h.UpdateSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Title", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Description", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Keywords", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGTitle", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGDescription", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OGImage", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "",
					ObjName:       "objname",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OBJName", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "",
					ObjPk:         "objpk",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing OBJPK", func(t *testing.T) {
			res, err := h.UpdateSEO(
				ctx, &pb.SEOMsg{
					Title:         "title",
					Description:   "description",
					Keywords:      "keywords",
					OGTitle:       "ogtitle",
					OGDescription: "ogdescription",
					OGImage:       "ogimage",
					ObjName:       "objname",
					ObjPk:         "",
				},
			)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)
}

func TestHandler_DeleteSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()
	name, pk := "name", "pk"

	t.Run(
		"Success", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			mockCtrl.EXPECT().
				DeleteSEO(gomock.Any(), name, pk).
				Return(nil).
				Times(1)

			res, err := h.DeleteSEO(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			mockCtrl.EXPECT().
				DeleteSEO(gomock.Any(), name, pk).
				Return(ctrl.ErrNotFound).
				Times(1)

			res, err := h.DeleteSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			req := &pb.GetSEOReq{Name: "name", Pk: "pk"}
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				DeleteSEO(gomock.Any(), name, pk).
				Return(newErr).
				Times(1)

			res, err := h.DeleteSEO(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing All", func(t *testing.T) {
			res, err := h.DeleteSEO(ctx, &pb.GetSEOReq{Name: "", Pk: ""})

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Name", func(t *testing.T) {
			res, err := h.DeleteSEO(ctx, &pb.GetSEOReq{Name: "", Pk: "pk"})

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing PK", func(t *testing.T) {
			res, err := h.DeleteSEO(ctx, &pb.GetSEOReq{Name: "name", Pk: ""})

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)
}
