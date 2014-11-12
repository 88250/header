// Copyright (c) 2014, B3log
//  
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//  
//     http://www.apache.org/licenses/LICENSE-2.0
//  
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"strings"
)

var pathSeparator = string(os.PathSeparator)

func AntPathMatch(pattern, path string) bool {
	if strings.HasPrefix(path, pathSeparator) != strings.HasPrefix(pattern, pathSeparator) {
		return false
	}

	pattDirs := strings.Split(pattern, pathSeparator)
	pathDirs := strings.Split(path, pathSeparator)

	pattIdxStart := 0
	pattIdxEnd := len(pattDirs) - 1
	pathIdxStart := 0
	pathIdxEnd := len(pathDirs) - 1

	// Match all elements up to the first **
	for pattIdxStart <= pattIdxEnd && pathIdxStart <= pathIdxEnd {
		patDir := pattDirs[pattIdxStart]

		if "**" == patDir {
			break
		}

		if !matchStrings(patDir, pathDirs[pathIdxStart]) {
			return false
		}
		pattIdxStart++
		pathIdxStart++
	}

	if pathIdxStart > pathIdxEnd {
		// Path is exhausted, only match if rest of pattern is * or **'s
		if pattIdxStart > pattIdxEnd {
			if strings.HasSuffix(pattern, pathSeparator) {
				return strings.HasSuffix(path, pathSeparator)
			} else {
				return !strings.HasSuffix(path, pathSeparator)
			}
		}
		if !false {
			return true
		}
		if pattIdxStart == pattIdxEnd && pattDirs[pattIdxStart] == "*" && strings.HasSuffix(path, pathSeparator) {
			return true
		}
		for i := pattIdxStart; i <= pattIdxEnd; i = i + 1 {
			if pattDirs[i] != "**" {
				return false
			}
		}
		return true
	} else if pattIdxStart > pattIdxEnd {
		// string not exhausted, but pattern is. Failure.
		return false
	} else if !false && "**" == pattDirs[pattIdxStart] {
		// Path start definitely matches due to "**" part in pattern.
		return true
	}

	// up to last '**'
	for pattIdxStart <= pattIdxEnd && pathIdxStart <= pathIdxEnd {
		patDir := pattDirs[pattIdxEnd]

		if patDir == "**" {
			break
		}
		if !matchStrings(patDir, pathDirs[pathIdxEnd]) {
			return false
		}
		pattIdxEnd--
		pathIdxEnd--
	}
	if pathIdxStart > pathIdxEnd {
		// string is exhausted
		for i := pattIdxStart; i <= pattIdxEnd; i = i + 1 {
			if pattDirs[i] != "**" {
				return false
			}
		}
		return true
	}

	for pattIdxStart != pattIdxEnd && pathIdxStart <= pathIdxEnd {
		patIdxTmp := -1

		for i := pattIdxStart + 1; i <= pattIdxEnd; i = i + 1 {
			if pattDirs[i] == "**" {
				patIdxTmp = i
				break
			}
		}
		if patIdxTmp == pattIdxStart+1 {
			// '**/**' situation, so skip one
			pattIdxStart++
			continue
		}
		// Find the pattern between padIdxStart & padIdxTmp in str between
		// strIdxStart & strIdxEnd
		patLength := (patIdxTmp - pattIdxStart - 1)
		strLength := (pathIdxEnd - pathIdxStart + 1)
		foundIdx := -1

	strLoop:
		for i := 0; i <= strLength-patLength; i = i + 1 {
			for j := 0; j < patLength; j = j + 1 {
				subPat := pattDirs[pattIdxStart+j+1]
				subStr := pathDirs[pathIdxStart+i+j]

				if !matchStrings(subPat, subStr) {
					continue strLoop
				}
			}

			foundIdx = pathIdxStart + i
			break
		}

		if foundIdx == -1 {
			return false
		}

		pattIdxStart = patIdxTmp
		pathIdxStart = foundIdx + patLength
	}

	for i := pattIdxStart; i <= pattIdxEnd; i = i + 1 {
		if pattDirs[i] != "**" {
			return false
		}
	}

	return true
}

func matchStrings(pattern, str string) bool {
	patArr := []byte(pattern)
	strArr := []byte(str)

	patIdxStart := 0
	patIdxEnd := len(patArr) - 1
	strIdxStart := 0
	strIdxEnd := len(strArr) - 1
	var ch byte

	containsStar := false

	for i := 0; i < len(patArr); i++ {
		if patArr[i] == '*' {
			containsStar = true
			break
		}
	}

	if !containsStar {
		// No '*'s, so we make a shortcut
		if patIdxEnd != strIdxEnd {
			return false // Pattern and string do not have the same size
		}
		for i := 0; i <= patIdxEnd; i++ {
			ch = patArr[i]
			if ch != '?' {
				if ch != strArr[i] {
					return false // Character mismatch
				}
			}
		}
		return true // string matches against pattern
	}

	if patIdxEnd == 0 {
		return true // Pattern contains only '*', which matches anything
	}

	// Process characters before first star
	for ch = patArr[patIdxStart]; ch != '*' && strIdxStart <= strIdxEnd; ch = patArr[patIdxStart] {
		if ch != '?' {
			if ch != strArr[strIdxStart] {
				return false // Character mismatch
			}
		}
		patIdxStart++
		strIdxStart++
	}
	if strIdxStart > strIdxEnd {
		// All characters in the string are used. Check if only '*'s are
		// left in the pattern. If so, we succeeded. Otherwise failure.
		for i := patIdxStart; i <= patIdxEnd; i = i + 1 {
			if patArr[i] != '*' {
				return false
			}
		}
		return true
	}

	// Process characters after last star
	for ch = patArr[patIdxEnd]; ch != '*' && strIdxStart <= strIdxEnd; ch = patArr[patIdxEnd] {
		if ch != '?' {
			if ch != strArr[strIdxEnd] {
				return false // Character mismatch
			}
		}
		patIdxEnd--
		strIdxEnd--
	}
	if strIdxStart > strIdxEnd {
		// All characters in the string are used. Check if only '*'s are
		// left in the pattern. If so, we succeeded. Otherwise failure.
		for i := patIdxStart; i <= patIdxEnd; i = i + 1 {
			if patArr[i] != '*' {
				return false
			}
		}
		return true
	}

	// process pattern between stars. padIdxStart and patIdxEnd point
	// always to a '*'.
	for patIdxStart != patIdxEnd && strIdxStart <= strIdxEnd {
		patIdxTmp := -1

		for i := patIdxStart + 1; i <= patIdxEnd; i = i + 1 {
			if patArr[i] == '*' {
				patIdxTmp = i
				break
			}
		}
		if patIdxTmp == patIdxStart+1 {
			// Two stars next to each other, skip the first one.
			patIdxStart++
			continue
		}
		// Find the pattern between padIdxStart & padIdxTmp in str between
		// strIdxStart & strIdxEnd
		patLength := (patIdxTmp - patIdxStart - 1)
		strLength := (strIdxEnd - strIdxStart + 1)
		foundIdx := -1

	strLoop:
		for i := 0; i <= strLength-patLength; i = i + 1 {
			for j := 0; j < patLength; j++ {
				ch = patArr[patIdxStart+j+1]
				if ch != '?' {
					if ch != strArr[strIdxStart+i+j] {
						continue strLoop
					}
				}
			}

			foundIdx = strIdxStart + i
			break
		}

		if foundIdx == -1 {
			return false
		}

		patIdxStart = patIdxTmp
		strIdxStart = foundIdx + patLength
	}

	// All characters in the string are used. Check if only '*'s are left
	// in the pattern. If so, we succeeded. Otherwise failure.
	for i := patIdxStart; i <= patIdxEnd; i = i + 1 {
		if patArr[i] != '*' {
			return false
		}
	}

	return true
}
