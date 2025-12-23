# Code Quality and Clarity Examples

## Code Commenting Examples

### Example 1: setTimeout vs setInterval
```javascript
// REASON: Using setTimeout instead of setInterval to prevent overlapping calls
// if the API response is slower than the interval
const scheduleNextCheck = () => {
  setTimeout(async () => {
    await checkSystemHealth();
    scheduleNextCheck();
  }, HEALTH_CHECK_INTERVAL);
};
```

### Example 2: Inventory Reservation
```javascript
// WHY: Business rule requires immediate inventory reservation to prevent overselling
// during the payment process window
async function reserveInventory(items, duration = 15 * 60 * 1000) {
  // NOTE: 15-minute default reservation window per business requirements
  const expiresAt = new Date(Date.now() + duration);
  return await InventoryService.reserve(items, expiresAt);
}
```

### Example 3: Date Validation
```javascript
// REASON: Manual validation needed because third-party library doesn't handle
// our specific date format requirements for international users
function validateDateFormat(dateString, locale) {
  // WHY: Different regions have different date format expectations
  const formats = getAcceptedFormatsForLocale(locale);
  return formats.some(format => moment(dateString, format, true).isValid());
}
```

## Safety Protocol Examples

### Wrong vs Right Approaches

❌ **WRONG**: Assume API endpoint structure
✅ **RIGHT**: "Could you confirm the expected API response format for the user endpoint?"

❌ **WRONG**: Delete seemingly unused code
✅ **RIGHT**: "I found this function that appears unused. Should I remove it or is it needed for future features?"

❌ **WRONG**: Overwrite existing configuration
✅ **RIGHT**: "I need to update the config file. Should I modify the existing settings or add new ones alongside?"

❌ **WRONG**: Guess at business logic
✅ **RIGHT**: "For the payment processing, should failed payments be retried automatically or require manual intervention?"
