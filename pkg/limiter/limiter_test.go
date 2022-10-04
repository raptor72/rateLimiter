package limiter

import (
	"strconv"
	"testing"
	"time"

	"github.com/raptor72/rateLimiter/config"
	"github.com/stretchr/testify/require"
)

func TestGetCountUnsetPattern(t *testing.T) {
	var expectedPatternCount int
	emptyPatterns := []string{"login", "password", "1234567", "192.168.1.1", "10.8.8.8"}

	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(c)

	for _, tc := range emptyPatterns {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			client.rdb.FlushAll(client.ctx)
			num, err := client.GetCountPattern(tc)
			require.NoError(t, err)
			require.Equal(t, expectedPatternCount, *num)
		})
	}
}

func TestGetCountSinglePattern(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "login", expected: 1},
		{input: "password", expected: 1},
		{input: "192.168.1.1", expected: 1},
	}

	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(c)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			client.rdb.FlushAll(client.ctx)
			_, err := client.rdb.Incr(client.ctx, tc.input+":").Result()
			if err != nil {
				t.Fatal(err)
			}
			num, err := client.GetCountPattern(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, *num)
		})
	}
}

func TestGetCountMultipleSameTimePattern(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "login", expected: 2},
		{input: "password", expected: 12},
		{input: "192.168.1.1", expected: 108},
	}

	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(c)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			client.rdb.FlushAll(client.ctx)
			client.rdb.IncrBy(client.ctx, tc.input+":", int64(tc.expected))
			num, err := client.GetCountPattern(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, *num)
		})
	}
}

func TestGetCountMultipleDifferentTimePattern(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "login", expected: 2},
		{input: "password", expected: 3},
		{input: "192.168.1.1", expected: 4},
	}
	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(c)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			client.rdb.FlushAll(client.ctx)
			for i := 0; i < tc.expected; i++ {
				time.Sleep(2 * time.Microsecond)
				now := time.Now()
				timestamp := now.Second()
				_, err := client.rdb.Incr(client.ctx, tc.input+":"+strconv.Itoa(timestamp)).Result()
				if err != nil {
					t.Fatal(err)
				}
			}
			num, err := client.GetCountPattern(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, *num)
		})
	}
}

func TestGetCountNotIntPattern(t *testing.T) {
	var expectedPatternCount int
	tests := []struct {
		input string
		value interface{}
	}{
		{input: "login", value: "efwe"},
		{input: "password", value: ""},
		{input: "192.168.1.1", value: "f43t$"},
	}

	c, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(c)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			client.rdb.FlushAll(client.ctx)
			client.rdb.Set(client.ctx, tc.input, tc.value, 2*time.Millisecond)
			num, err := client.GetCountPattern(tc.input)
			require.NoError(t, err)
			require.Equal(t, expectedPatternCount, *num)
		})
	}
}
