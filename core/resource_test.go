package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestResourceQuantity_Add(t *testing.T) {
	tests := []struct {
		name     string
		base     core.ResourceQuantity
		other    core.ResourceQuantity
		expected core.ResourceQuantity
	}{
		{
			name: "正常な加算",
			base: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			other: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: core.ResourceQuantity{
				Money: 150, Food: 75, Wood: 45, Iron: 30, Mana: 15,
			},
		},
		{
			name: "ゼロとの加算",
			base: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			other: core.ResourceQuantity{},
			expected: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
		},
		{
			name: "負の値との加算",
			base: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			other: core.ResourceQuantity{
				Money: -30, Food: -10, Wood: -5, Iron: -15, Mana: -3,
			},
			expected: core.ResourceQuantity{
				Money: 70, Food: 40, Wood: 25, Iron: 5, Mana: 7,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Add(tt.other)
			if result != tt.expected {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResourceQuantity_Sub(t *testing.T) {
	tests := []struct {
		name     string
		base     core.ResourceQuantity
		other    core.ResourceQuantity
		expected core.ResourceQuantity
	}{
		{
			name: "正常な減算",
			base: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			other: core.ResourceQuantity{
				Money: 30, Food: 20, Wood: 10, Iron: 5, Mana: 3,
			},
			expected: core.ResourceQuantity{
				Money: 70, Food: 30, Wood: 20, Iron: 15, Mana: 7,
			},
		},
		{
			name: "ゼロとの減算",
			base: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			other: core.ResourceQuantity{},
			expected: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
		},
		{
			name: "負の値になる減算",
			base: core.ResourceQuantity{
				Money: 50, Food: 30, Wood: 20, Iron: 10, Mana: 5,
			},
			other: core.ResourceQuantity{
				Money: 80, Food: 40, Wood: 30, Iron: 15, Mana: 10,
			},
			expected: core.ResourceQuantity{
				Money: -30, Food: -10, Wood: -10, Iron: -5, Mana: -5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Sub(tt.other)
			if result != tt.expected {
				t.Errorf("Sub() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResourceQuantity_CanPurchase(t *testing.T) {
	tests := []struct {
		name     string
		treasury core.ResourceQuantity
		price    core.ResourceQuantity
		expected bool
	}{
		{
			name: "十分な資源がある場合",
			treasury: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			price: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: true,
		},
		{
			name: "ちょうど同じ資源量の場合",
			treasury: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			price: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			expected: true,
		},
		{
			name: "Moneyが不足している場合",
			treasury: core.ResourceQuantity{
				Money: 40, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			price: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: false,
		},
		{
			name: "Foodが不足している場合",
			treasury: core.ResourceQuantity{
				Money: 100, Food: 20, Wood: 30, Iron: 20, Mana: 10,
			},
			price: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: false,
		},
		{
			name: "複数リソースが不足している場合",
			treasury: core.ResourceQuantity{
				Money: 40, Food: 20, Wood: 10, Iron: 5, Mana: 2,
			},
			price: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: false,
		},
		{
			name: "価格が0の場合",
			treasury: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			price:    core.ResourceQuantity{},
			expected: true,
		},
		{
			name:     "国庫が空の場合、価格が0なら購入可能",
			treasury: core.ResourceQuantity{},
			price:    core.ResourceQuantity{},
			expected: true,
		},
		{
			name:     "国庫が空の場合、何か価格があれば購入不可",
			treasury: core.ResourceQuantity{},
			price: core.ResourceQuantity{
				Money: 1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.treasury.CanPurchase(tt.price)
			if result != tt.expected {
				t.Errorf("CanPurchase() = %v, want %v", result, tt.expected)
			}
		})
	}
}
