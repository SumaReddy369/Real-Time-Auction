# Integration test script - requires Docker (PostgreSQL + Redis) to be running
# Run: docker-compose up -d
# Then: .\scripts\integration_test.ps1

$baseUrl = "http://localhost:8081"

Write-Host "=== Real-Time Auction Integration Tests ===" -ForegroundColor Cyan

# Health check
Write-Host "`n1. Health check..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
    Write-Host "   OK" -ForegroundColor Green
} catch {
    Write-Host "   FAIL: $_" -ForegroundColor Red
    exit 1
}

# Create auction
Write-Host "`n2. Create auction..." -ForegroundColor Yellow
$createBody = @{
    title = "Vintage Camera"
    description = "Rare 1960s camera"
    start_price = 50.00
    duration_min = 5
} | ConvertTo-Json

try {
    $auction = Invoke-RestMethod -Uri "$baseUrl/api/v1/auctions" -Method Post -Body $createBody -ContentType "application/json"
    Write-Host "   Created auction ID: $($auction.id)" -ForegroundColor Green
    $auctionId = $auction.id
} catch {
    Write-Host "   FAIL: $_" -ForegroundColor Red
    exit 1
}

# Get auction
Write-Host "`n3. Get auction..." -ForegroundColor Yellow
try {
    $get = Invoke-RestMethod -Uri "$baseUrl/api/v1/auctions/$auctionId" -Method Get
    if ($get.id -eq $auctionId -and $get.title -eq "Vintage Camera") {
        Write-Host "   OK" -ForegroundColor Green
    } else {
        Write-Host "   FAIL: unexpected response" -ForegroundColor Red
    }
} catch {
    Write-Host "   FAIL: $_" -ForegroundColor Red
}

# Place bid
Write-Host "`n4. Place bid..." -ForegroundColor Yellow
$bidBody = @{
    bidder_id = "user123"
    amount = 75.00
} | ConvertTo-Json

try {
    $bid = Invoke-RestMethod -Uri "$baseUrl/api/v1/auctions/$auctionId/bids" -Method Post -Body $bidBody -ContentType "application/json"
    Write-Host "   Bid placed: $($bid.amount) by $($bid.bidder_id)" -ForegroundColor Green
} catch {
    Write-Host "   FAIL: $_" -ForegroundColor Red
}

# List bids
Write-Host "`n5. List bids..." -ForegroundColor Yellow
try {
    $bids = Invoke-RestMethod -Uri "$baseUrl/api/v1/auctions/$auctionId/bids" -Method Get
    $count = if ($bids -is [array]) { $bids.Count } else { 1 }
    if ($count -ge 1) {
        Write-Host "   OK - $count bid(s)" -ForegroundColor Green
    } else {
        Write-Host "   FAIL: expected at least 1 bid" -ForegroundColor Red
    }
} catch {
    Write-Host "   FAIL: $_" -ForegroundColor Red
}

# Invalid bid (lower than current)
Write-Host "`n6. Reject invalid bid (too low)..." -ForegroundColor Yellow
$lowBidBody = @{
    bidder_id = "user456"
    amount = 50.00
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri "$baseUrl/api/v1/auctions/$auctionId/bids" -Method Post -Body $lowBidBody -ContentType "application/json"
    Write-Host "   FAIL: should have rejected low bid" -ForegroundColor Red
} catch {
    Write-Host "   OK - correctly rejected" -ForegroundColor Green
}

Write-Host "`n=== All integration tests completed ===" -ForegroundColor Cyan
