#!/bin/bash

echo "🚀 Loading demo data into LocalConnect database..."
echo ""

# Check if PostgreSQL is running
if ! docker ps | grep -q localconnect-postgres; then
    echo "❌ Error: PostgreSQL container is not running!"
    echo "   Start it with: docker compose up -d"
    exit 1
fi

echo "📦 Loading seed data..."
docker exec -i localconnect-postgres psql -U postgres -d localconnect < migrations/003_seed_data.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Demo data loaded successfully!"
    echo ""
    echo "📊 What was created:"
    echo "   • 4 Demo clients"
    echo "   • 12 Demo workers with profiles"
    echo "   • Categories: Carpenter, Painter, Labour, Cook, Home Tuition, Teacher, Househelp"
    echo "   • Locations: New York, Los Angeles, Chicago, Houston, Phoenix, and more"
    echo "   • Multiple reviews and ratings for each worker"
    echo "   • Upvotes from clients"
    echo "   • Some contact requests"
    echo ""
    echo "🔑 Demo login credentials (all passwords: 'password'):"
    echo "   Workers:"
    echo "     - robert@example.com (Carpenter - New York)"
    echo "     - jennifer@example.com (Painter - Los Angeles)"
    echo "     - maria@example.com (Cook - Houston)"
    echo ""
    echo "   Clients:"
    echo "     - sarah@example.com"
    echo "     - michael@example.com"
    echo ""
else
    echo ""
    echo "❌ Error loading demo data!"
    exit 1
fi
