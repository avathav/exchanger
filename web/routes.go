package web

import "exchanger/web/handlers"

func (s *Server) Routes() {
	s.r.GET("/rates", handlers.GetRatesHandler(s.exchangeService))
	s.r.GET("/exchange", handlers.GetExchangeHandler(s.storage))
}
