package mapper

import (
	"github.com/JMURv/seo/api/grpc/v1/gen"
	md "github.com/JMURv/seo/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PagesToProto(req []*md.Page) []*gen.PageMsg {
	var res []*gen.PageMsg
	for _, v := range req {
		res = append(res, PageToProto(v))
	}
	return res
}

func PageToProto(req *md.Page) *gen.PageMsg {
	return &gen.PageMsg{
		Slug:      req.Slug,
		Title:     req.Title,
		Href:      req.Href,
		CreatedAt: timestamppb.New(req.CreatedAt),
		UpdatedAt: timestamppb.New(req.UpdatedAt),
	}
}

func ProtoToPage(req *gen.PageMsg) *md.Page {
	return &md.Page{
		Slug:      req.Slug,
		Title:     req.Title,
		Href:      req.Href,
		CreatedAt: req.CreatedAt.AsTime(),
		UpdatedAt: req.UpdatedAt.AsTime(),
	}
}
