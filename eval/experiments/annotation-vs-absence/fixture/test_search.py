from app import ITEMS, search


def test_search_matches_case_insensitively():
    assert [i["id"] for i in search("meeting")] == [4]


def test_search_no_match_is_empty():
    assert search("gardening") == []


def test_search_empty_query_returns_all():
    assert search("") == ITEMS
