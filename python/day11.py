from collections import defaultdict

import utils


def convert(n: int) -> list[int]:
    if n == 0:
        return [1]

    s = str(n)
    if len(s) % 2 == 1:
        return [n * 2024]

    return [int(s[: len(s) // 2]), int(s[len(s) // 2 :])]


def read() -> list[int]:
    with utils.read(day=11, sample=False) as f:
        return list(map(int, f.read().strip().split()))


def solve(times: int) -> None:
    nums = read()
    count: dict[int, int] = defaultdict(lambda: 0)
    for n in nums:
        count[n] += 1

    for _ in range(times):
        new_count = count.copy()
        for k, c in count.items():
            new_count[k] -= c
            for v in convert(k):
                new_count[v] += c
        count = new_count

    s = 0
    for c in count.values():
        s += c
    print(s)


def part1() -> None:
    solve(25)


def part2() -> None:
    solve(75)


if __name__ == "__main__":
    part1()
    part2()
