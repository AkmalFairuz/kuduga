package service

import (
	"errors"
	"github.com/akmalfairuz/finance/module/checker"
	"github.com/akmalfairuz/finance/module/digiflazz"
	"github.com/akmalfairuz/finance/server/repository"
	"strings"
)

type CheckerService struct {
	digiflazzClient   *digiflazz.Client
	checkerRepository *repository.CheckerRepository
	handlers          map[string]CheckerFunc
}

func NewCheckerService(checkerRepository *repository.CheckerRepository) *CheckerService {
	service := &CheckerService{handlers: map[string]CheckerFunc{}, checkerRepository: checkerRepository}
	service.RegisterHandler("pln", service.checkPLN)
	service.RegisterHandler("dana", service.checkDana)
	service.RegisterHandler("gopay", service.checkGopay)
	service.RegisterHandler("ovo", service.checkOvo)
	service.RegisterHandler("shopeepay", service.checkShopeePay)
	return service
}

type CheckerFunc func(input map[string]string) (CheckerResult, error)

func (s *CheckerService) Check(checkerId string, input map[string]string) (CheckerResult, error) {
	handler, ok := s.handlers[checkerId]
	if !ok {
		return nil, errors.New("checker handler not found")
	}
	cachedResult, err := s.checkerRepository.Get(checkerId, input)
	if err != nil {
		return nil, err
	}
	if cachedResult != nil {
		return cachedResult, nil
	}
	result, err := handler(input)
	if err != nil {
		return nil, err
	}
	if err := s.checkerRepository.Save(checkerId, input, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CheckerService) RegisterHandlerWithTtl(checkerId string, handler CheckerFunc, ttl int64) {
	s.handlers[checkerId] = handler
}

func (s *CheckerService) RegisterHandler(checkerId string, handler CheckerFunc) {
	s.RegisterHandlerWithTtl(checkerId, handler, 0)
}

func (s *CheckerService) checkPLN(input map[string]string) (CheckerResult, error) {
	resp, err := s.digiflazzClient.InquiryPLN(input["id"])
	if err != nil {
		return nil, err
	}
	result := CheckerResult{}
	result.Add("Nama", resp.Name)
	parts := strings.Split(resp.SegmentPower, "/")
	if len(parts) == 2 {
		result.Add("Tarif", parts[0])
		daya := strings.TrimLeft(parts[1], "0") + " W"
		result.Add("Daya", daya)
	} else {
		result.Add("Tarif & Daya", resp.SegmentPower)
	}
	return result, nil
}

func (s *CheckerService) checkDana(input map[string]string) (CheckerResult, error) {
	resp, err := checker.CheckDana(input["no"])
	if err != nil {
		return nil, err
	}
	resp = strings.TrimLeft(resp, "DNID ")
	return NewCheckerResultSingle("Nama", resp), nil
}

func (s *CheckerService) checkGopay(input map[string]string) (CheckerResult, error) {
	resp, err := checker.CheckGopay(input["no"])
	if err != nil {
		return nil, err
	}
	return NewCheckerResultSingle("Nama", resp), nil
}

func (s *CheckerService) checkShopeePay(input map[string]string) (CheckerResult, error) {
	resp, err := checker.CheckShopeepay(input["no"])
	if err != nil {
		return nil, err
	}
	return NewCheckerResultSingle("Nama", resp), nil
}

func (s *CheckerService) checkOvo(input map[string]string) (CheckerResult, error) {
	resp, err := checker.CheckOvo(input["no"])
	if err != nil {
		return nil, err
	}
	return NewCheckerResultSingle("Nama", resp), nil
}

func NewCheckerResultSingle(key string, val string) CheckerResult {
	ret := CheckerResult{}
	ret.Add(key, val)
	return ret
}

type CheckerResult [][]string

func (r *CheckerResult) Add(key string, val string) {
	*r = append(*r, []string{key, val})
}
