import dataclasses
import heapq
import itertools
from collections import defaultdict

import utils

DOT = "."


def read_layout() -> list[int | str]:
    layout: list[int | str] = []
    with utils.read(day=9, sample=False) as f:
        blocks = f.read().strip()
        c = 0
        for i, b in enumerate(blocks):
            if i % 2 == 1:
                layout.extend(itertools.repeat(DOT, int(b)))
            else:
                layout.extend(itertools.repeat(c, int(b)))
                c += 1
    return layout


Dots = dict[int, list[int]]


@dataclasses.dataclass
class Num:
    num: int
    at: int
    n: int


def read_layout2() -> tuple[Dots, list[Num]]:
    dots = defaultdict(list)
    nums = []
    with utils.read(day=9, sample=False) as f:
        blocks = f.read().strip()
        c = 0
        pos = 0
        for i, b in enumerate(blocks):
            if i % 2 == 1:
                dots[int(b)].append(pos)
            else:
                nums.append(Num(c, pos, int(b)))
                c += 1
            pos += int(b)

    for k in dots:
        heapq.heapify(dots[k])

    return dots, nums


def part1() -> None:
    layout = read_layout()

    i = 0
    j = len(layout) - 1

    ans = 0
    while i < j:
        while i < j and layout[j] == DOT:
            j -= 1
        while i < j and layout[i] != DOT:
            i += 1

        c = layout[j]
        layout[j] = "."
        layout[i] = c

        i += 1
        j -= 1

    for i, c in enumerate(layout):
        if isinstance(c, int):
            ans += i * c

    print(ans)


def part2() -> None:
    dots, nums = read_layout2()

    ans = 0
    for num in nums[::-1]:
        key: int | None = None
        min_at: int | None = None
        for k in range(num.n, 10):
            if len(dots[k]):
                at = heapq.heappop(dots[k])
                heapq.heappush(dots[k], at)

                if (min_at is None or min_at > at) and (at < num.at):
                    min_at = at
                    key = k

        if key is not None:
            # found
            at = heapq.heappop(dots[key])

            dots_size = key - num.n
            dots_at = at + num.n
            if dots_size > 0:
                heapq.heappush(dots[dots_size], dots_at)

            num.at = at

        ans += num.num * (num.at * num.n + num.n * (num.n - 1) // 2)

    print(ans)


if __name__ == "__main__":
    part1()
    part2()
