package services

import "github.com/go-resty/resty/v2"

type IRestService interface {
	Get(url string) (*resty.Response, error)
	Post(url string, body interface{}) (*resty.Response, error)
	Put(url string, body interface{}) (*resty.Response, error)
	Delete(url string) (*resty.Response, error)
	Patch(url string, body interface{}) (*resty.Response, error)
	Head(url string) (*resty.Response, error)
}

type RestyRestService struct {
	client resty.Client
}

// TODO wrap resty in unspecific struct
func NewRestyRestService() *RestyRestService {
	return &RestyRestService{
		client: *resty.New(),
	}
}

func (s *RestyRestService) Get(url string) (*resty.Response, error) {
	return s.client.R().Get(url)
}

func (s *RestyRestService) Post(url string, body interface{}) (*resty.Response, error) {
	return s.client.R().SetBody(body).Post(url)
}

func (s *RestyRestService) Put(url string, body interface{}) (*resty.Response, error) {
	return s.client.R().SetBody(body).Put(url)
}

func (s *RestyRestService) Delete(url string) (*resty.Response, error) {
	return s.client.R().Delete(url)
}

func (s *RestyRestService) Patch(url string, body interface{}) (*resty.Response, error) {
	return s.client.R().SetBody(body).Patch(url)
}

func (s *RestyRestService) Head(url string) (*resty.Response, error) {
	return s.client.R().Head(url)
}
