package biz

import "errors"

// ErrNotFound adalah error standar saat resource tidak ditemukan.
var ErrNotFound = errors.New("resource not found")
