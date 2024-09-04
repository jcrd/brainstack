export function parseTaskText(text) {
    let tags = []
    let i = 0
    let tags_idx = []
    let last_text_idx

    const words = text.split(" ")

    words.forEach((word) => {
        if (word.startsWith("#")) {
            tags.push(word.slice(1))
            tags_idx.push(i)
        } else {
            last_text_idx = i
        }
        i++
    })

    const parsed = {
        text: text,
        parsed_text: words.filter((_, i) => !(i > last_text_idx)).join(" "),
        tags: tags.filter((_, i) => tags_idx[i] > last_text_idx),
    }

    return parsed
}
