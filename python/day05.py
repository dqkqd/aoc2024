import dataclasses
import itertools
from collections import defaultdict

import utils


@dataclasses.dataclass
class Rules:
    rules: dict[int, set[int]] = dataclasses.field(
        default_factory=lambda: defaultdict(set)
    )

    def add_rule(self, line: str) -> None:
        lhs, rhs, *_ = map(int, line.split("|"))
        self.rules[lhs].add(rhs)

    def topology_sort(self, subset: set[int]) -> dict[int, int]:
        visited: set[int] = set()
        answer: list[int] = []

        def dfs(v: int) -> None:
            visited.add(v)
            for u in self.rules[v] & subset:
                if u not in visited:
                    dfs(u)
            answer.append(v)

        for v in subset:
            if v not in visited:
                dfs(v)

        answer.reverse()

        return {v: order for (order, v) in enumerate(answer)}


@dataclasses.dataclass
class Page:
    data: list[int]

    @classmethod
    def from_str(cls, line: str) -> "Page":
        page = map(int, line.split(","))
        return Page(list(page))

    @property
    def subset(self) -> set[int]:
        return set(self.data)


def read_rules_and_pages() -> tuple[Rules, list[Page]]:
    rules = Rules()
    pages: list[Page] = []

    read_rule = True
    with utils.read(day=5, sample=False) as reader:
        for line in reader.readlines():
            if line.strip() == "":
                read_rule = False
                continue

            if read_rule:
                rules.add_rule(line)
            else:
                pages.append(Page.from_str(line))

    return rules, pages


def is_sorted(data: list[int]) -> bool:
    return all(lhs <= rhs for lhs, rhs in itertools.pairwise(data))


def part1() -> None:
    rules, pages = read_rules_and_pages()
    ans = 0
    for page in pages:
        order = rules.topology_sort(page.subset)
        ordered_page = [order[p] for p in page.data]
        if is_sorted(ordered_page):
            ans += page.data[len(page.data) // 2]
    print(ans)


def part2() -> None:
    rules, pages = read_rules_and_pages()
    ans = 0
    for page in pages:
        order = rules.topology_sort(page.subset)
        ordered_page = [order[p] for p in page.data]
        if not is_sorted(ordered_page):
            ordered_page = sorted(ordered_page)
            middle_page = next(
                (
                    k
                    for k, v in order.items()
                    if v == ordered_page[len(ordered_page) // 2]
                ),
                0,
            )
            ans += middle_page

    print(ans)


if __name__ == "__main__":
    part1()
    part2()
