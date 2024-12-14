import functools
import itertools
import string

import utils

MapAtenna = list[str]
Point = tuple[int, int]


def get_atennas(m: MapAtenna, a: str) -> list[Point]:
    return [(x, y) for x, line in enumerate(m) for y, c in enumerate(line) if c == a]


def possible_antennas() -> str:
    return string.ascii_letters + string.digits


def good_point(m: MapAtenna, p: Point) -> bool:
    return p[0] >= 0 and p[0] < len(m) and p[1] >= 0 and p[1] < len(m[0])


def read_map() -> MapAtenna:
    with utils.read(day=8, sample=False) as f:
        map_attenas = []
        for line in f.readlines():
            if len(line.strip()) == 0:
                continue
            atennas = line.strip()
            map_attenas.append(atennas)
    return map_attenas


def part1() -> None:
    m = read_map()
    s: set[Point] = set()
    for a in possible_antennas():
        points = get_atennas(m, a)
        for p1, p2 in itertools.product(points, points):
            if p1 == p2:
                continue
            px = (p2[0] * 2 - p1[0], p2[1] * 2 - p1[1])
            py = (p1[0] * 2 - p2[0], p1[1] * 2 - p2[1])
            s.add(px)
            s.add(py)

    good = functools.partial(good_point, m)
    ans = len(list(filter(good, s)))
    print(ans)


def good_point2(p: Point, atennas: list[Point]) -> bool:
    for p1, p2 in itertools.product(atennas, atennas):
        if p1 == p2:
            continue

        if p in (p1, p2):
            return True

        x1, y1 = p[0] - p1[0], p[1] - p1[1]
        x2, y2 = (p[0] - p2[0], p[1] - p2[1])

        if x1 == 0:
            x1, y1 = y1, x1
            x2, y2 = y2, x2

        v = x2 / x1
        if x1 * v == x2 and y1 * v == y2:
            return True

    return False


def part2() -> None:
    m = read_map()
    ans = 0
    all_atennas = [get_atennas(m, a) for a in possible_antennas()]
    for x in range(len(m)):
        for y in range(len(m[0])):
            p = (x, y)
            good = functools.partial(good_point2, p)
            if any(map(good, all_atennas)):
                ans += 1
    print(ans)


if __name__ == "__main__":
    part1()
    part2()
