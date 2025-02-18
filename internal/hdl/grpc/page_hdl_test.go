package grpc

import (
	"context"
	"errors"
	pb "github.com/JMURv/seo/api/grpc/v1/gen"
	"github.com/JMURv/seo/internal/ctrl"
	model "github.com/JMURv/seo/internal/models"
	utils "github.com/JMURv/seo/internal/models/mapper"
	"github.com/JMURv/seo/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestHandler_ListPages(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	var expected []*model.Page
	ctx := context.Background()
	req := &pb.EmptySEO{}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().ListPages(gomock.Any()).Return(expected, nil).Times(1)

			res, err := h.ListPages(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"Nil req", func(t *testing.T) {
			res, err := h.ListPages(ctx, nil)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().ListPages(gomock.Any()).Return(expected, newErr).Times(1)

			res, err := h.ListPages(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)
}

func TestHandler_GetPage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()
	slug := "slug"
	expected := &model.Page{}
	req := &pb.SlugSEO{Slug: "slug"}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().GetPage(gomock.Any(), slug).Return(expected, nil).Times(1)

			res, err := h.GetPage(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"InvalidArgument", func(t *testing.T) {
			req.Slug = ""
			res, err := h.GetPage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			req.Slug = "slug"
			mockCtrl.EXPECT().GetPage(gomock.Any(), slug).Return(expected, ctrl.ErrNotFound).Times(1)

			res, err := h.GetPage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().GetPage(gomock.Any(), slug).Return(expected, newErr).Times(1)

			res, err := h.GetPage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)
}

func TestHandler_CreatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()

	req := &pb.PageMsg{
		Slug:  "slug",
		Title: "name",
		Href:  "href",
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), utils.ProtoToPage(req)).
				Return("slug", nil).
				Times(1)

			res, err := h.CreatePage(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"Nil req", func(t *testing.T) {
			res, err := h.CreatePage(ctx, nil)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), utils.ProtoToPage(req)).
				Return("", ctrl.ErrNotFound).
				Times(1)

			res, err := h.CreatePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), utils.ProtoToPage(req)).
				Return("", newErr).
				Times(1)

			res, err := h.CreatePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Slug", func(t *testing.T) {
			req.Slug = ""
			res, err := h.CreatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Slug = "slug"
		},
	)

	t.Run(
		"InvalidArgument - Missing Title", func(t *testing.T) {
			req.Title = ""
			res, err := h.CreatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Title = "name"
		},
	)

	t.Run(
		"InvalidArgument - Missing Href", func(t *testing.T) {
			req.Href = ""
			res, err := h.CreatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Href = "href"
		},
	)

}

func TestHandler_UpdatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()

	slug := "slug"
	req := &pb.PageWithSlugMsg{
		Slug: slug,
		Page: &pb.PageMsg{
			Slug:  "slug",
			Title: "name",
			Href:  "href",
		},
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, utils.ProtoToPage(req.Page)).
				Return(nil).
				Times(1)

			res, err := h.UpdatePage(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"Nil req", func(t *testing.T) {
			res, err := h.UpdatePage(ctx, nil)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, utils.ProtoToPage(req.Page)).
				Return(ctrl.ErrNotFound).
				Times(1)

			res, err := h.UpdatePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, utils.ProtoToPage(req.Page)).
				Return(newErr).
				Times(1)

			res, err := h.UpdatePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Slug", func(t *testing.T) {
			req.Slug = ""
			res, err := h.UpdatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Slug = "slug"
		},
	)

	t.Run(
		"InvalidArgument - Missing Title", func(t *testing.T) {
			req.Page.Title = ""
			res, err := h.UpdatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Page.Title = "name"
		},
	)

	t.Run(
		"InvalidArgument - Missing Href", func(t *testing.T) {
			req.Page.Href = ""
			res, err := h.UpdatePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Page.Href = "href"
		},
	)

}

func TestHandler_DeletePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New("", mockCtrl, ssoCtrl)

	ctx := context.Background()
	slug := "slug"
	req := &pb.SlugSEO{Slug: slug}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(nil).
				Times(1)

			res, err := h.DeletePage(ctx, req)
			assert.Nil(t, err)
			assert.NotNil(t, res)
		},
	)

	t.Run(
		"Nil req", func(t *testing.T) {
			res, err := h.DeletePage(ctx, nil)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(ctrl.ErrNotFound).
				Times(1)

			res, err := h.DeletePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.NotFound, status.Code(err))
		},
	)

	t.Run(
		"Internal Error", func(t *testing.T) {
			newErr := errors.New("new error")
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(newErr).
				Times(1)

			res, err := h.DeletePage(ctx, req)
			assert.Nil(t, res)
			assert.Equal(t, codes.Internal, status.Code(err))
		},
	)

	t.Run(
		"InvalidArgument - Missing Slug", func(t *testing.T) {
			req.Slug = ""
			res, err := h.DeletePage(ctx, req)

			assert.Nil(t, res)
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			req.Slug = "slug"
		},
	)

}
