import dataclasses
from collections import defaultdict
from functools import cached_property

import utils


@dataclasses.dataclass(unsafe_hash=True)
class Point:
    x: int
    y: int


@dataclasses.dataclass
class Map:
    data: list[list[int]]

    @cached_property
    def h(self) -> int:
        return len(self.data)

    @cached_property
    def w(self) -> int:
        return len(self.data[0])

    def zeros(self) -> list[Point]:
        return [
            Point(x, y)
            for x, line in enumerate(self.data)
            for y, c in enumerate(line)
            if c == 0
        ]

    def adjacents(self, point: Point) -> list[Point]:
        def good(p: Point) -> bool:
            return p.x >= 0 and p.x < self.h and p.y >= 0 and p.y < self.w

        return list(
            filter(
                good,
                [
                    Point(point.x, point.y - 1),
                    Point(point.x, point.y + 1),
                    Point(point.x - 1, point.y),
                    Point(point.x + 1, point.y),
                ],
            )
        )

    def count_nine(self, point: Point) -> int:
        assert self.data[point.x][point.y] == 0

        currents = {point}
        for score in range(1, 10):
            new_currents = set()
            for p in currents:
                for adj in self.adjacents(p):
                    if self.data[adj.x][adj.y] == score:
                        new_currents.add(adj)
            currents = new_currents

        return len(currents)

    def count_nine2(self, point: Point) -> int:
        assert self.data[point.x][point.y] == 0

        score_map: dict[Point, int] = defaultdict(lambda: 0)
        score_map[point] = 1

        currents = {point}
        for score in range(1, 10):
            new_currents = set()
            for p in currents:
                for adj in self.adjacents(p):
                    if self.data[adj.x][adj.y] == score:
                        score_map[adj] += score_map[p]
                        new_currents.add(adj)
            currents = new_currents

        s = 0
        for p in currents:
            s += score_map[p]
        return s


def read_map() -> Map:
    with utils.read(day=10, sample=False) as f:
        data = [
            list(map(int, line.strip())) for line in f.readlines() if len(line.strip())
        ]
    return Map(data)


def part1() -> None:
    m = read_map()

    scores = 0
    for zero in m.zeros():
        scores += m.count_nine(zero)

    print(scores)


def part2() -> None:
    m = read_map()

    scores = 0
    for zero in m.zeros():
        scores += m.count_nine2(zero)

    print(scores)


if __name__ == "__main__":
    part1()
    part2()
