package request

type ArticleRequest struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type GetDetailArticleRequest struct {
	Id int `param:"id" validate:"required"`
}
