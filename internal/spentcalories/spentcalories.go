package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректный формат данных: ожидается 3 элемента, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть положительным")
	}

	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, errors.New("вид активности не указан")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования длительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть положительной")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	dist := height * stepLengthCoefficient
	dist = (dist * float64(steps)) / mInKm
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var dist, speed, calories float64
	switch activity {
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		cals, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчёта калорий для бега: %w", err)
		}
		calories = cals
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		cals, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчёта калорий для ходьбы: %w", err)
		}
		calories = cals
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}

// ниже для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 || height <= 0 {
		return 0, errors.New("вес и рост должны быть положительными числами")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность тренировки должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	// Формула: (вес * средняя скорость * длительность в минутах) / минут в часе
	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// ниже для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 || height <= 0 {
		return 0, errors.New("вес и рост должны быть положительными числами")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность тренировки должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationMinutes) / minInH
	calories := baseCalories * walkingCaloriesCoefficient
	return calories, nil
}
