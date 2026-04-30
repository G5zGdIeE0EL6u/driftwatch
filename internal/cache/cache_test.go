package cache

import (
	"testing"
	"time"
)

func TestSet_Get_Hit(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("key", "value")
	v, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v.(string) != "value" {
		t.Fatalf("expected 'value', got %v", v)
	}
}

func TestGet_Miss(t *testing.T) {
	c := New(5 * time.Second)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestSet_Expired(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("key", "value")
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestDelete(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("key", 42)
	c.Delete("key")
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestFlush_RemovesExpired(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("a", 1)
	c.SetWithTTL("b", 2, 10*time.Second)
	time.Sleep(5 * time.Millisecond)
	c.Flush()
	if c.Len() != 1 {
		t.Fatalf("expected 1 item after flush, got %d", c.Len())
	}
	_, ok := c.Get("b")
	if !ok {
		t.Fatal("expected 'b' to still be present")
	}
}

func TestSetWithTTL_Override(t *testing.T) {
	c := New(5 * time.Second)
	c.SetWithTTL("key", "short", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected short-TTL entry to be expired")
	}
}
