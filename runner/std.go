package runner

const ARRAYS = `
let map = fn(arr, func) {
	let mut result = []
	for i in 0..len(arr) {
		result = append(result, func(arr[i]))
	}
	return result
};

let forEach = fn(arr, func) {
	for i in 0..len(arr) {
		func(arr[i])
	}
};

let filter = fn(arr, func) {
	let mut result = []
	for i in 0..len(arr) {
		if func(arr[i]) {
			result = append(result, arr[i]);
		}
	}
	return result
};

let reduce = fn(arr, func, mut acc) {
	for i in 0..len(arr) {
		acc = func(acc, arr[i])
	}
	return acc
};

let find = fn(arr, func) {
	for i in 0..len(arr) {
		if func(arr[i]) {
			return arr[i]
		}
	}
};

let findIndex = fn(arr, func) {
	for i in 0..len(arr) {
		if func(arr[i]) {
			return i
		}
	}
	return -1
};

let some = fn(arr, func) {
	for i in 0..len(arr) {
		if func(arr[i]) {
			return true
		}
	}
	return false
};

let every = fn(arr, func) {
	for i in 0..len(arr) {
		if !func(arr[i]) {
			return false
		}
	}
	return true
};
`

const MATH = `
let min = fn(a, b) {
	if a < b {
		return a
	}
	return b
};

let max = fn(a, b) {
	if a > b {
		return a
	}
	return b
};

let abs = fn(x) {
	if x < 0 {
		return -x
	}
	return x
};

let clamp = fn(x, min, max) {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
};

let pow = fn(x, n) {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	let mut result = 1
	for i in 0..(n-1) {
		result = result * x
	}
	return result
};
`

const STD = ARRAYS + MATH
