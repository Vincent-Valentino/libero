# Project Improvement Plan: League Page UI & Structure

**Date:** 2025-03-31

**Goal:** Refine the UI/UX of the League Page (`src/leagues/LeaguePage.vue` and its child components) and prepare the component structure for easier integration of real data fetching in the future.

## Phase 1: UI/UX Refinement

1.  **Analyze Current State:** Review the existing layout (masonry grid), spacing, alignment, card design, and overall aesthetic based on the code and visual inspection.
2.  **Implement UI Changes (using Tailwind CSS):**
    *   **Consistency:** Ensure consistent padding, margins, font sizes, and card styles across `UpcomingMatches`, `PlayerShowdown`, `LeagueTable`, and `LeagueHistory` components.
    *   **Visual Hierarchy:** Improve the visual distinction between sections and within cards (e.g., clearer headings, better use of whitespace).
    *   **`LeaguePage.vue` Layout:**
        *   Review and potentially adjust the `.masonry-container` styling (gap, columns, item styling).
        *   Consider alternative layouts (simple grid/flex) if masonry proves difficult or doesn't achieve the desired look.
    *   **Component-Specific Styling:** Refine styles within each child component for better presentation (e.g., table styling, card appearance).
    *   **Responsiveness:** Verify and improve layout adaptation to various screen sizes.

## Phase 2: Structural Preparation for Future Data Fetching

1.  **Maintain Clear Props:** Ensure well-defined TypeScript types for all data props passed from `LeaguePage` to child components.
2.  **Implement Conditional Rendering Placeholders:** Structure each child component to handle different states by adding placeholder elements/logic:
    *   **Loading State:** Add `v-if="isLoading"` blocks for future skeleton loaders or "Loading..." text.
    *   **Empty State:** Ensure existing "No data available" messages (`v-if="!data || data.length === 0"`) are styled appropriately.
    *   **Error State:** Add `v-if="error"` blocks for future error message displays.
3.  **(Future Suggestion) Data Service Layer:** Plan to create a separate service file (e.g., `src/services/leagueApi.ts`) later to encapsulate API fetching logic, improving separation of concerns in `LeaguePage.vue`.

## Implementation Order

1.  Focus on UI refinements (Phase 1).
2.  Integrate structural preparations (Phase 2) alongside UI changes where appropriate.
3.  Defer actual data fetching implementation and the data service layer creation.

## Next Steps

*   Switch to 'code' mode to begin implementing the UI refinements.