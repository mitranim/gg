package gg_test

import (
	"testing"

	"github.com/mitranim/gg"
	"github.com/mitranim/gg/gtest"
)

func TestChanInit(t *testing.T) {
	defer gtest.Catch(t)

	gg.ChanInit((*chan string)(nil))

	var tar chan string
	gtest.Eq(gg.ChanInit(&tar), tar)

	gtest.NotZero(tar)
	gtest.Eq(cap(tar), 0)

	prev := tar
	gtest.Eq(gg.ChanInit(&tar), prev)
	gtest.Eq(tar, prev)
}

func TestChanInitCap(t *testing.T) {
	defer gtest.Catch(t)

	gg.ChanInitCap((*chan string)(nil), 1)

	var tar chan string
	gtest.Eq(gg.ChanInitCap(&tar, 3), tar)

	gtest.NotZero(tar)
	gtest.Eq(cap(tar), 3)

	prev := tar
	gtest.Eq(gg.ChanInitCap(&tar, 5), prev)
	gtest.Eq(tar, prev)
	gtest.Eq(cap(prev), 3)
	gtest.Eq(cap(tar), 3)
}

func TestSendOpt(t *testing.T) {
	defer gtest.Catch(t)

	var tar chan string
	gg.SendOpt(tar, `one`)
	gg.SendOpt(tar, `two`)
	gg.SendOpt(tar, `three`)

	tar = make(chan string, 1)
	gg.SendOpt(tar, `one`)
	gg.SendOpt(tar, `two`)
	gg.SendOpt(tar, `three`)

	gtest.Eq(<-tar, `one`)
}

func TestSendZeroOpt(t *testing.T) {
	defer gtest.Catch(t)

	var tar chan string
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)

	tar = make(chan string, 1)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)
	gg.SendZeroOpt(tar)

	val, ok := <-tar
	gtest.Zero(val)
	gtest.True(ok)
}
