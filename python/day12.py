import dataclasses
import enum
import itertools
from functools import cached_property

import utils


class Side(enum.IntEnum):
    Horizontal = 0
    Vertical = 1


@dataclasses.dataclass(order=True)
class Edge:
    side: Side
    x: int
    y: int


@dataclasses.dataclass(unsafe_hash=True)
class Point:
    x: int
    y: int

    def shared_edge(self, point: "Point") -> Edge | None:
        if self.x == point.x:
            if self.y == point.y + 1:
                return Edge(Side.Vertical, self.x, self.y)
            if point.y == self.y + 1:
                return Edge(Side.Vertical, self.x, point.y)

        if self.y == point.y:
            if self.x == point.x + 1:
                return Edge(Side.Horizontal, self.x, self.y)
            if point.x == self.x + 1:
                return Edge(Side.Horizontal, point.x, self.y)

        return None


@dataclasses.dataclass
class Map:
    data: list[str]

    @cached_property
    def h(self) -> int:
        return len(self.data)

    @cached_property
    def w(self) -> int:
        return len(self.data[0])

    def inside(self, p: Point) -> bool:
        return p.x >= 0 and p.x < self.h and p.y >= 0 and p.y < self.w

    def adjacents(self, point: Point) -> list[Point]:
        return [
            Point(point.x, point.y - 1),
            Point(point.x, point.y + 1),
            Point(point.x - 1, point.y),
            Point(point.x + 1, point.y),
        ]

    def same(self, lhs: Point, rhs: Point) -> bool:
        if not self.inside(lhs) or not self.inside(rhs):
            return False
        return self.data[lhs.x][lhs.y] == self.data[rhs.x][rhs.y]

    def connectable(self, lhs: Edge, rhs: Edge) -> bool:
        if lhs.side != rhs.side:
            return False

        match lhs.side:
            case Side.Horizontal if lhs.x == rhs.x and abs(lhs.y - rhs.y) <= 1:
                p1, p2 = Point(lhs.x, lhs.y), Point(rhs.x, rhs.y)
                p3, p4 = Point(lhs.x - 1, lhs.y), Point(rhs.x - 1, rhs.y)

                return (
                    (not self.inside(p1) and not self.inside(p2) and self.same(p3, p4))
                    or (
                        not self.inside(p3)
                        and not self.inside(p4)
                        and self.same(p1, p2)
                    )
                    or (self.same(p1, p2) or self.same(p3, p4))
                )

            case Side.Vertical if lhs.y == rhs.y and abs(lhs.x - rhs.x) <= 1:
                p1, p2 = Point(lhs.x, lhs.y), Point(rhs.x, rhs.y)
                p3, p4 = Point(lhs.x, lhs.y - 1), Point(rhs.x, rhs.y - 1)

                return (
                    (not self.inside(p1) and not self.inside(p2) and self.same(p3, p4))
                    or (
                        not self.inside(p3)
                        and not self.inside(p4)
                        and self.same(p1, p2)
                    )
                    or (self.same(p1, p2) or self.same(p3, p4))
                )

        return False

    def visit(self, point: Point) -> set[Point]:
        visited: set[Point] = set()
        visited.add(point)

        def dfs(p: Point) -> None:
            for adj in filter(self.inside, self.adjacents(p)):
                if adj in visited or not self.same(point, adj):
                    continue
                visited.add(adj)
                dfs(adj)

        dfs(point)

        return visited


def read_map() -> Map:
    with utils.read(day=12, sample=False) as f:
        data = [line.strip() for line in f.readlines() if len(line.strip()) > 0]
    return Map(data)


def part1() -> None:
    m = read_map()
    visited: set[Point] = set()

    ans = 0

    for x, y in itertools.product(range(m.h), range(m.w)):
        point = Point(x, y)
        if point in visited:
            continue

        local_visited = m.visit(point)

        area = len(local_visited)
        perimeter = 0

        for p in local_visited:
            for adj in m.adjacents(p):
                if not m.same(p, adj):
                    perimeter += 1

        ans += area * perimeter

        visited = visited | local_visited

    print(ans)


def part2() -> None:
    m = read_map()
    visited: set[Point] = set()

    ans = 0

    for x, y in itertools.product(range(m.h), range(m.w)):
        point = Point(x, y)
        if point in visited:
            continue

        local_visited = m.visit(point)

        area = len(local_visited)

        edges: list[Edge] = []

        for p in local_visited:
            for adj in m.adjacents(p):
                if not m.same(p, adj):
                    edge = p.shared_edge(adj)
                    assert edge is not None
                    edges.append(edge)

        edges = sorted(edges)

        connected_edges: list[list[Edge]] = []

        for e in edges:
            found = False
            for cedges in connected_edges:
                if any(m.connectable(e, ce) for ce in cedges):
                    found = True
                    cedges.append(e)
            if not found:
                connected_edges.append([e])

        ans += area * len(connected_edges)

        visited = visited | local_visited

    print(ans)


if __name__ == "__main__":
    part1()
    part2()
