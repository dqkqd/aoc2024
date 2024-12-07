import bisect
import contextlib
import copy
import dataclasses
import enum
from collections import defaultdict
from collections.abc import Generator
from functools import cached_property

import utils


class MapItem(enum.Enum):
    Guard = 0
    Obstacle = 1
    Nothing = 2

    @classmethod
    def from_str(cls, s: str) -> "MapItem":
        if s == "^":
            return cls.Guard
        if s == "#":
            return cls.Obstacle
        return cls.Nothing


class Direction(enum.Enum):
    Up = 0
    Down = 1
    Left = 2
    Right = 3

    def turn(self) -> "Direction":
        match self:
            case Direction.Up:
                return Direction.Right
            case Direction.Down:
                return Direction.Left
            case Direction.Left:
                return Direction.Up
            case Direction.Right:
                return Direction.Down


@dataclasses.dataclass(unsafe_hash=True)
class Position:
    x: int
    y: int
    direction: Direction

    def move(self) -> "Position":
        position = copy.copy(self)

        match self.direction:
            case Direction.Up:
                position.x -= 1
            case Direction.Down:
                position.x += 1
            case Direction.Left:
                position.y -= 1
            case Direction.Right:
                position.y += 1

        return position

    def turn(self) -> "Position":
        position = copy.copy(self)
        position.direction = position.direction.turn()
        return position

    def out_of_bound(self, map_items: "Map") -> bool:
        return (
            self.x < 0
            or self.x >= len(map_items.data)
            or self.y < 0
            or self.y >= len(map_items.data[0])
        )

    def valid(self, map_items: "Map") -> bool:
        return (
            not self.out_of_bound(map_items)
            and map_items.data[self.x][self.y] is not MapItem.Obstacle
        )

    def farthest_pos(self, map_items: "Map") -> "Position | None":
        assert self.valid(map_items)
        position = copy.copy(self)

        match self.direction:
            case Direction.Up:
                obstacle_index = (
                    bisect.bisect(map_items.vertical_obstacles[position.y], position.x)
                    - 1
                )
                if obstacle_index < 0:
                    return None
                position.x = (
                    map_items.vertical_obstacles[position.y][obstacle_index] + 1
                )
            case Direction.Down:
                obstacle_index = bisect.bisect(
                    map_items.vertical_obstacles[position.y], position.x
                )
                if obstacle_index >= len(map_items.vertical_obstacles[position.y]):
                    return None
                position.x = (
                    map_items.vertical_obstacles[position.y][obstacle_index] - 1
                )
            case Direction.Left:
                obstacle_index = (
                    bisect.bisect(
                        map_items.horizontal_obstacles[position.x], position.y
                    )
                    - 1
                )
                if obstacle_index < 0:
                    return None
                position.y = (
                    map_items.horizontal_obstacles[position.x][obstacle_index] + 1
                )
            case Direction.Right:
                obstacle_index = bisect.bisect(
                    map_items.horizontal_obstacles[position.x], position.y
                )
                if obstacle_index >= len(map_items.horizontal_obstacles[position.x]):
                    return None
                position.y = (
                    map_items.horizontal_obstacles[position.x][obstacle_index] - 1
                )

        return position


@dataclasses.dataclass
class Map:
    data: list[list[MapItem]] = dataclasses.field(default_factory=list)
    vertical_obstacles: dict[int, list[int]] = dataclasses.field(
        default_factory=lambda: defaultdict(list)
    )
    horizontal_obstacles: dict[int, list[int]] = dataclasses.field(
        default_factory=lambda: defaultdict(list)
    )

    def add_line(self, line: str) -> None:
        line = line.strip()
        if len(line) == 0:
            return

        x = len(self.data)
        map_items = list(map(MapItem.from_str, line))

        self.data.append(map_items)
        for y, item in enumerate(map_items):
            if item is MapItem.Obstacle:
                bisect.insort_left(self.vertical_obstacles[y], x)
                bisect.insort_left(self.horizontal_obstacles[x], y)

    @cached_property
    def guard_position(self) -> Position:
        for x, items in enumerate(self.data):
            for y, item in enumerate(items):
                if item == MapItem.Guard:
                    return Position(x=x, y=y, direction=Direction.Up)
        raise ValueError

    def run(self) -> set[tuple[int, int]]:
        pos = self.guard_position

        positions = set()

        while not pos.out_of_bound(self):
            positions.add((pos.x, pos.y))
            next_pos = pos.move()
            if next_pos.out_of_bound(self):
                break
            pos = next_pos if next_pos.valid(self) else pos.turn()

        return positions

    @contextlib.contextmanager
    def with_obstacle(self, x: int, y: int) -> Generator[None]:
        assert self.data[x][y] is not MapItem.Obstacle
        assert x not in self.vertical_obstacles[y]
        assert y not in self.horizontal_obstacles[x]

        self.data[x][y] = MapItem.Obstacle
        bisect.insort_left(self.vertical_obstacles[y], x)
        bisect.insort_left(self.horizontal_obstacles[x], y)

        yield

        self.data[x][y] = MapItem.Nothing
        self.vertical_obstacles[y].remove(x)
        self.horizontal_obstacles[x].remove(y)

    def is_loop(self) -> bool:
        pos = self.guard_position
        traces: set[Position] = set()
        while True:
            if pos in traces:
                return True
            traces.add(pos)
            farthest_pos = pos.farthest_pos(self)
            if farthest_pos is None:
                return False
            pos = farthest_pos.turn()


def read_map() -> Map:
    items = Map()
    with utils.read(day=6, sample=False) as f:
        for line in f.readlines():
            items.add_line(line)
    return items


def part1() -> None:
    items = read_map()
    positions = items.run()
    print(len(positions))


def part2() -> None:
    items = read_map()
    positions = items.run()
    ans = 0
    for x, y in positions:
        if x == items.guard_position.x and y == items.guard_position.y:
            continue
        with items.with_obstacle(x, y):
            if items.is_loop():
                ans += 1
    print(ans)


if __name__ == "__main__":
    part1()
    part2()
