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
		if (func(arr[i])) {
			result = append(result, arr[i]);
		}
	}
	return result
};
`
