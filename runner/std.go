package runner

const STD = `
let map = fn(arr, func) {
	let mut result = []
	for i in 0..len(arr) {
		result = append(result, func(arr[i]))
	}
	return result
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
`
