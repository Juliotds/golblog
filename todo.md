# TODO
- [ ] Explain the origin of Sayou Gauma
- [X] Reorder so the newest post comes in the top
- [X] show date of the post in the home page
- [ ] Add the logo in the page
- [ ] Add versioning to the blog
- [ ] Add update for lambda functions

## UI Improvements
- [ ] Add post excerpts/descriptions in the post list so readers get a preview before clicking
- [ ] Make tags clickable to filter posts by tag
- [ ] Add social links (GitHub, etc.) to the hero section or footer
- [ ] Add a brief bio/description in the hero section — "Writing about projects." is too generic
- [ ] Add a "View all posts" link if the home page only shows recent posts
- [ ] Add a search bar or search functionality for posts
- [ ] Hero section has too much empty space — consider adding an avatar or profile image

## Color Scheme & Typography
- [X] Replace hardcoded hex values with CSS custom properties (variables) so the palette is easy to change in one place
- [X] Add a more vibrant accent/highlight color for interactive elements (buttons, active links, CTAs) — the current teal (#18608C) is too muted
- [X] Improve tag contrast — tag backgrounds (#000a2e on #011140) are very hard to distinguish from the page background
- [X] Use a dedicated web font for body text in blog posts (e.g. a readable serif like Lora or a clean sans-serif like Inter) — the current system font stack is generic and the monospace headers make body text feel inconsistent
- [X] Load JetBrains Mono from a CDN (Google Fonts or Bunny Fonts) — it's currently referenced by name only and silently falls back to Courier New if not installed on the visitor's machine
- [X] Audit color contrast ratios for accessibility — light blue (#8ec5d8) and medium blue (#6aaac0) text on dark navy may fail WCAG AA for smaller text sizes
- [X] Establish a clearer typographic hierarchy — heading sizes (h1–h4) are all monospace but the size jumps between them could be more deliberate to guide the reader's eye