# NOX Racket Scraper (Go)

A simple Go scraper that collects NOX padel racket data from noxsport.com and outputs it to JSON files.

## What it does
- Visits a category/series page (e.g. CLassic, Advanced, ...)
- Extracts for each racket:
    - `brand`, `model`, `price`, `imageUrl`, `racketPage`
- Opens each racket detail page and parses:
    - `weight`, `shape`, `material` (from the product description)
- Writes results to:
    - `NOXRackets<Series>.json` (e.g. `NOXRacketsSerieUltralight.json`)

## Output format
Each item looks like:
```json
  {
    "brand": "NOX",
    "model": "Equation SOFT Advanced",
    "price": "124,99 €",
    "imageUrl": "https://noxsport.com/cdn/shop/files/equation-soft-advanced-pequsadv26-8435778902751-4653383.png",
    "racketPage": "https://noxsport.com/en/collections/serie-advance/products/pala-equation-soft-advanced",
    "weight": "360-375g",
    "shape": "Round",
    "material": "Fiber Glass",
    "series": "SerieAdvance"
  }