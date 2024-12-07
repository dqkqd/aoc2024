import itertools
from collections.abc import Generator

from utils import read


def safe(levels: list[int]) -> bool:
    windows = list(itertools.pairwise(levels))
    is_diff = all(1 <= abs(a - b) <= 3 for (a, b) in windows)
    is_incr = all(a < b for (a, b) in windows)
    is_decr = all(a > b for (a, b) in windows)
    return is_diff and (is_incr or is_decr)


def safe_rm(levels: list[int]) -> bool:
    return safe(levels) or any(
        safe(levels[:i] + levels[i + 1 :]) for i in range(len(levels))
    )


def read_input() -> Generator[list[int]]:
    with read(day=2, sample=False) as f:
        reports = (line.split() for line in f.readlines())
    return (list(map(int, levels)) for levels in reports)


def part1() -> None:
    safe_cnt = sum(map(safe, read_input()))
    print(safe_cnt)


def part2() -> None:
    safe_cnt = sum(map(safe_rm, read_input()))
    print(safe_cnt)


if __name__ == "__main__":
    part1()
    part2()
