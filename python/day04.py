import itertools

from utils import read


def chains(x: int, y: int) -> list[list[tuple[int, int]]]:
    return [
        [(x + p * dx, y + p * dy) for p in range(4)]
        for dx in [-1, 0, 1]
        for dy in [-1, 0, 1]
    ]


def count_xmas(words: list[str], x: int, y: int) -> int:
    cs = chains(x, y)

    h = len(words)
    w = len(words[0])

    cnt = 0
    for c in cs:
        if any(x < 0 or x >= h for x, _ in c) or any(y < 0 or y >= w for _, y in c):
            continue
        word = "".join(words[x][y] for x, y in c)
        cnt += word == "XMAS"

    return cnt


def getmap(words: list[str], x: int, y: int) -> list[str] | None:
    h = len(words)
    w = len(words[0])
    if x + 2 >= h:
        return None
    if y + 2 >= w:
        return None
    return [
        words[x][y] + words[x][y + 1] + words[x][y + 2],
        words[x + 1][y] + words[x + 1][y + 1] + words[x + 1][y + 2],
        words[x + 2][y] + words[x + 2][y + 1] + words[x + 2][y + 2],
    ]


def mapgood(word: list[str]) -> bool:
    if word[1][1] != "A":
        return False
    pos = [(0, 0), (0, 2), (2, 2), (2, 0)]
    pos = list(itertools.chain(pos, pos))
    gen = ["".join(word[x][y] for x, y in pos[off : off + 4]) for off in range(4)]
    return any(w == "MMSS" for w in gen)


def part1():
    with read(day=4, sample=False) as f:
        words = [line.strip() for line in f.readlines()]
    h = len(words)
    w = len(words[0])
    ans = sum(count_xmas(words, x, y) for x in range(h) for y in range(w))
    print(ans)


def part2():
    with read(day=4, sample=False) as f:
        words = [line.strip() for line in f.readlines()]
    h = len(words)
    w = len(words[0])
    ans = 0
    for x in range(h):
        for y in range(w):
            word = getmap(words, x, y)
            if word is None:
                continue
            ans += mapgood(word)
    print(ans)


if __name__ == "__main__":
    part1()
    part2()
