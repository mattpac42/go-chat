# Mobile-Friendly Project Rename Design

## Problem Statement
Current inline edit uses Enter/Escape keyboard shortcuts, which aren't discoverable on mobile. Users may not know how to confirm title changes on touch devices.

**Constraints:**
- Compact project card (limited horizontal space)
- Mobile-first priority
- Desktop should also support (but can keep keyboard shortcuts)
- Trash icon exists in bottom right
- MVP needs to be simple

---

## Design Solution: Hybrid Approach

### Wireframe Description

```
EDIT MODE - Mobile & Desktop Layout:

┌─────────────────────────────┐
│ [×] Project Name Input    │  ← Input field spans full width
└─────────────────────────────┘
│ [✓ Save] [✕ Cancel]       │  ← Action buttons (mobile visible)
│         [🗑 Delete]        │  ← Trash icon repositioned below
└─────────────────────────────┘

EDIT MODE - Desktop (Alternative - Compact):
┌─────────────────────────────┐
│ [×] Project Name Input    │
│ [✓] [✕]           [🗑]    │  ← Icons only (with hover labels)
└─────────────────────────────┘
```

### Recommended Approach: Context-Aware Display

**Mobile (< 768px):**
- Show explicit "Save" and "Cancel" buttons below input
- Use text labels (not icons) for clarity
- Stack vertically or place side-by-side (test both)
- Make buttons 44px+ height for touch targets
- Move trash icon below action buttons (separate concern)

**Desktop (≥ 768px):**
- Show compact icon buttons: ✓ (checkmark) and ✕ (X)
- Keep keyboard shortcuts (Enter/Escape) as primary
- Icons have hover tooltips: "Save (Enter)" and "Cancel (Escape)"
- Trash icon stays visible or tucked in overflow menu

---

## Implementation Details

### Mobile Implementation (Recommended)

```
Layout: Flex column, gap: 8px

Input row:
- Full width input field
- Close/cancel [×] button on right (tap to cancel quickly)

Action buttons row:
- [Save] button (primary color, 48px height)
- [Cancel] button (secondary color, 48px height)
- Equal width, gap: 8px

Trash icon:
- Below buttons, full width, danger color
- Or: Move to separate "More options" menu
```

### Desktop Implementation

```
Layout: Inline with icons

Input field | [✓] [✕] [🗑]
- Buttons use 24px icons
- On hover: background highlight + tooltip
- Keyboard still primary (Enter/Escape)
```

---

## Design Decisions & Rationale

| Decision | Rationale |
|----------|-----------|
| **Text buttons on mobile** | Icons alone are less discoverable; text removes ambiguity for unfamiliar users |
| **44px+ button height** | WCAG touch target minimum; reduces misclicks on mobile |
| **Separate trash from confirm** | Different action (destructive) deserves distinct placement; prevents accidental deletes |
| **Keep keyboard on desktop** | Faster for power users; reduces visual clutter in compact layouts |
| **Close [×] in input** | Quick cancel option; users expect this pattern (modals, search bars) |
| **Checkmark/X icons** | Universal symbols; more compact than labels for desktop |

---

## Responsive Behavior

**Breakpoint: 768px**

- Below 768px: Full-width text buttons, vertical stack
- Above 768px: Icon buttons, inline layout
- Button labels fade to icons at 768px
- Trash icon repositioned at 768px threshold

---

## Accessibility Considerations

1. **Button labels**: Text buttons have clear labels (not just icons)
2. **Touch targets**: Min 44px x 44px for mobile
3. **Keyboard support**: Tab navigation, Enter/Escape shortcuts
4. **Visual focus**: Clear focus indicators on all buttons
5. **Color**: Don't rely on color alone; use icons + text on mobile
6. **ARIA labels**: Save button has aria-label="Save project name"

---

## Code Structure Example

```jsx
{isEditing ? (
  <div className="edit-mode">
    {/* Input field with close button */}
    <div className="input-wrapper">
      <input
        value={name}
        onChange={setName}
        onKeyDown={handleKeyDown} // Enter/Escape
        placeholder="Project name"
      />
      <button
        onClick={cancelEdit}
        aria-label="Cancel"
        className="close-btn"
      >
        ×
      </button>
    </div>

    {/* Mobile: Text buttons | Desktop: Icon buttons */}
    <div className="actions">
      <button
        onClick={saveName}
        className="btn btn-primary btn-save"
      >
        Save
      </button>
      <button
        onClick={cancelEdit}
        className="btn btn-secondary btn-cancel"
      >
        Cancel
      </button>
    </div>

    {/* Trash: Below actions or in menu */}
    <button
      onClick={deleteProject}
      className="btn btn-danger btn-delete"
      aria-label="Delete project"
    >
      Delete
    </button>
  </div>
) : null}
```

---

## CSS Breakpoint Strategy

```css
/* Mobile (default) */
.edit-mode {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn {
  flex: 1;
  height: 48px;
  font-size: 16px; /* Prevent zoom on iOS */
}

/* Desktop */
@media (min-width: 768px) {
  .actions {
    flex-direction: row;
    gap: 4px;
  }

  .btn {
    width: auto;
    padding: 6px 12px;
    height: 32px;
    font-size: 14px;
  }

  .btn::before {
    content: attr(data-icon); /* Show icons instead of text */
  }
}
```

---

## MVP Recommendation

**Start with: Explicit Text Buttons on Mobile**

1. Add "Save" and "Cancel" buttons below input (mobile)
2. Keep keyboard shortcuts (Enter/Escape) for both
3. For desktop: show same buttons initially (simplest)
4. Optional polish: Switch to icon-only buttons on desktop at 768px+

**Why this MVP approach:**
- Clear, discoverable on mobile
- Works on all screen sizes
- No progressive enhancement complexity
- Easy to refine later (add icons, reorganize trash)
- Tests well with users

---

## Testing Recommendations

1. **Mobile usability test**: Can new users find and use Save/Cancel?
2. **Desktop test**: Do keyboard shortcuts work as expected?
3. **Touch test**: Button hit targets (44px minimum)
4. **Accessibility audit**: Keyboard nav, screen reader compatibility
5. **Compare**: Text vs. icon buttons on mobile (A/B test if time allows)

---

## Next Steps

1. Implement MVP with text buttons (mobile & desktop)
2. Test with 3-5 mobile users
3. Gather feedback on button placement
4. Optional: Add responsive icon buttons for desktop at 768px+
5. Consider: Move trash icon to secondary menu if space is tight
