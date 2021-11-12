// export function debounce (func: (e: any) => void, timeout = 300): () => any {
//   let timer
//   return (...args) => {
//     clearTimeout(timer)
//     timer = setTimeout(() => { func.apply(this, args) }, timeout)
//   }
// }

/**
 * fuzzyStringMatch performs fuzzy string matching
 * retrieved from: https://www.willmcgugan.com/blog/tech/post/sublime-text-like-fuzzy-matching-in-javascript/
 * @param {string} text text string
 * @param {string} searchText text substring to match
 * @returns {string} matched text
 */
export function fuzzyStringMatch(text: string, searchText: string): string {
	// remove spaces, lower case the search so the search is case insensitive
	const search = searchText.replace(/ /g, '').toLowerCase()
	const tokens = []
	let pos = 0

	// go through each character in the text
	for (let n = 0; n < text.length; n++) {
		const textChar = text[n]
		if (pos < search.length && textChar.toLowerCase() == search[pos]) {
			pos++
		}
		tokens.push(textChar)
	}

	// if are characters remaining in the search text, return an empty string to indicate no match
	if (pos != search.length) {
		return ''
	}
	return tokens.join('')
}

export class Storage {
	static getItem(key: string): string | null {
		return localStorage.getItem(key)
	}

	static setItem(key: string, value: string): void {
		localStorage.setItem(key, value)
	}
}
