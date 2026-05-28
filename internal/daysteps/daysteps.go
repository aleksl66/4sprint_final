package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"4sprint_final/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат данных: ожидается 2 элемента, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования длительности: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	step, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка разбора данных:", err)
		return ""
	}
	if step <= 0 {
		return ""
	}

	distance := float64(step) * stepLength
	distance = distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(step, weight, height, duration)
	if err != nil {
		log.Println("Ошибка вычисления калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		step,
		distance,
		calories,
	)
}
