package index

/*
Intersect returns all postings that are present in all lists
*/
func InterSect(posts ...PostingsList) PostingsList {
	results := posts[0]
	for len(posts) > 1 && len(results) != 0 {
		posts = posts[1:]
		results = InterSectTwo(results, posts[0])
	}
	return results
}

func InterSectTwo(first, second PostingsList) PostingsList {
	var result []DocID
	for len(first) > 0 && len(second) > 0 {
		if first[0] == second[0] {
			result = append(result, first[0])
			first = first[1:]
			second = second[1:]
		} else {
			if first[0].String() < second[0].String() {
				first = first[1:]
			} else {
				second = second[1:]
			}
		}
	}
	return result
}
