"""A tiny bookmarks store: add items, search titles by text."""

ITEMS = [
    {"id": 1, "title": "Team retro notes", "tags": ["work", "meetings"]},
    {"id": 2, "title": "Quarterly planning doc", "tags": ["work", "planning"]},
    {"id": 3, "title": "Sourdough starter guide", "tags": ["cooking"]},
    {"id": 4, "title": "Meeting-free Fridays proposal", "tags": ["work", "meetings", "planning"]},
    {"id": 5, "title": "Weeknight pasta ideas", "tags": ["cooking", "quick"]},
    {"id": 6, "title": "1:1 agenda template", "tags": ["meetings"]},
]


def add_item(title, tags=None):
    item = {"id": max(i["id"] for i in ITEMS) + 1, "title": title, "tags": list(tags or [])}
    ITEMS.append(item)
    return item


def search(query):
    """Case-insensitive title search."""
    q = query.lower()
    return [i for i in ITEMS if q in i["title"].lower()]
