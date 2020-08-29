package biz

import "context"

/*
站外
@api_doc_tags: outside
@api_doc_http_method: GET
@api_doc_relative_paths: /zlexample/outside/hi
*/
func Outside(ctx context.Context) (resp string) {
	return "this is outside"
}
