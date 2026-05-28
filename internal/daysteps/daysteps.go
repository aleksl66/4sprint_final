package daysteps

import (
	"fmt"
	"internal/spentcalories/spentcalories.go"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректный формат данных: ожидается 2 элемента, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным %d", steps)
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
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
		fmt.Println("Ошибка разбора данных:", err)
		return ""
	}
	if step <= 0 {
		return ""
	}
	distance := stepLength * float64(step)
	distance = distance / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка вычисления калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

}
