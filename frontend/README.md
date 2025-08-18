# WB Tech Frontend

A simple, modern web interface for looking up orders by ID in the WB Tech system.

## Features

- **Simple Order Lookup**: Enter an order ID to view order details
- **Clean Interface**: Modern, responsive design
- **Order Details**: Comprehensive order information display
- **Error Handling**: User-friendly error messages
- **Responsive Design**: Works on desktop, tablet, and mobile

## Technology Stack

- **HTML5**: Semantic markup
- **CSS3**: Modern styling with gradients and animations
- **JavaScript (ES6+)**: Vanilla JS with async/await
- **Font Awesome**: Icons
- **Google Fonts**: Inter font family

## File Structure

```
frontend/
├── index.html          # Main HTML interface
├── styles.css          # CSS styles
├── script.js           # JavaScript functionality
└── README.md           # This file
```

## Getting Started

### Local Development

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd wbtech-go/frontend
   ```

2. **Open in browser**

   ```bash
   # Using Python (if available)
   python3 -m http.server 3000

   # Using Node.js (if available)
   npx serve -s . -l 3000

   # Or simply open index.html in your browser
   ```

3. **Access the application**
   - Open `http://localhost:3000` in your browser

### Usage

1. **Enter Order ID**: Type an order ID (e.g., b1, b2, b3...) in the input field
2. **Lookup Order**: Click "Lookup Order" or press Enter
3. **View Details**: Order information will be displayed below
4. **Close View**: Click the X button to hide the order details

## API Integration

The frontend communicates with the backend through:

- `GET /order/:order_uid` - Get specific order details

## Features in Detail

### Order Lookup

- Simple input field for order ID
- Enter key support for quick lookup
- Loading indicator during API calls

### Order Display

- **Order Information**: ID, track number, date, customer ID
- **Payment Details**: Amount, currency, delivery cost
- **Delivery Information**: Customer name, city, address
- **Items List**: Product details with prices and status

### Error Handling

- Clear error messages for invalid IDs
- Network error notifications
- Auto-hiding error messages

## Responsive Design

The frontend is fully responsive and works on:

- **Desktop**: Full layout with all features
- **Tablet**: Optimized for medium screens
- **Mobile**: Mobile-first design with touch-friendly controls

## Browser Support

- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

## Development

### Adding New Features

1. **HTML**: Add new elements to `index.html`
2. **CSS**: Style new elements in `styles.css`
3. **JavaScript**: Add functionality in `script.js`

### Styling Guidelines

- Use CSS custom properties for consistent theming
- Follow modern CSS practices with Flexbox and Grid
- Implement mobile-first responsive design
- Use smooth transitions and hover effects

### JavaScript Guidelines

- Use ES6+ features
- Follow functional programming principles
- Implement proper error handling
- Use async/await for API calls

## Local Testing

### Without Docker

1. **Start your backend service** (make sure it's running on port 8080)
2. **Open the frontend** in your browser
3. **Test with real order IDs** from your system

### With Docker Backend

1. **Start the backend service**:

   ```bash
   docker-compose up backend
   ```

2. **Open the frontend** in your browser
3. **Test order lookup** functionality

## Troubleshooting

### Common Issues

1. **API Connection Errors**

   - Check if backend is running on port 8080
   - Verify the `/order/:id` endpoint works
   - Check browser console for errors

2. **Styling Issues**

   - Clear browser cache
   - Check CSS file loading
   - Verify Font Awesome CDN

3. **JavaScript Errors**
   - Check browser console
   - Verify API responses
   - Check for syntax errors

### Debug Mode

Enable debug logging by opening browser console and checking:

- Network requests to `/order/:id`
- JavaScript errors
- Console logs

## Contributing

1. Follow the existing code style
2. Test on multiple browsers
3. Ensure responsive design works
4. Update documentation as needed

## License

This frontend is part of the WB Tech project and follows the same license terms.
