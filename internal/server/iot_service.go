package server

import (
	"context"

	"github.com/go-redis/redis/v8"

	tempPb "NakrayIoT/proto"
)

type TemperatureService struct {
	tempPb.UnimplementedTemperatureServiceServer
	redisClient *redis.Client
}

func NewTemperatureService(redisClient *redis.Client) *TemperatureService {
	return &TemperatureService{redisClient: redisClient}
}

func (s *TemperatureService) RecordTemperature(ctx context.Context, req *tempPb.RecordTemperatureRequest) (*tempPb.RecordTemperatureResponse, error) {
	key := "temperature:" + req.GetSensorId()
	temperature := req.GetTemperature()

	err := s.redisClient.Set(ctx, key, temperature, 0).Err()
	if err != nil {
		return nil, err
	}

	err = s.redisClient.SAdd(ctx, "sensors", req.GetSensorId()).Err()
	if err != nil {
		return nil, err
	}

	return &tempPb.RecordTemperatureResponse{Success: true}, nil
}

func (s *TemperatureService) GetTemperature(ctx context.Context, req *tempPb.GetTemperatureRequest) (*tempPb.GetTemperatureResponse, error) {
	key := "temperature:" + req.GetSensorId()

	temperature, err := s.redisClient.Get(ctx, key).Float64()
	if err == redis.Nil {
		return &tempPb.GetTemperatureResponse{Found: false}, nil
	} else if err != nil {
		return nil, err
	}
	return &tempPb.GetTemperatureResponse{Temperature: temperature, Found: true}, nil
}

func (s *TemperatureService) GetAllTemperatures(ctx context.Context, req *tempPb.GetAllTemperaturesRequest) (*tempPb.GetAllTemperaturesResponse, error) {
	sensorIDs, err := s.redisClient.SMembers(ctx, "sensors").Result()
	if err != nil {
		return nil, err
	}

	var sensors []*tempPb.SensorTemperature
	for _, sensorID := range sensorIDs {
		key := "temperature:" + sensorID

		temperature, err := s.redisClient.Get(ctx, key).Float64()
		if err == redis.Nil {
			continue // Если температура для датчика отсутствует, пропускаем.
		} else if err != nil {
			return nil, err
		}

		sensors = append(sensors, &tempPb.SensorTemperature{
			SensorId:    sensorID,
			Temperature: temperature,
		})
	}

	return &tempPb.GetAllTemperaturesResponse{Sensors: sensors}, nil
}
