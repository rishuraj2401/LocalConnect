# 🧪 Testing Guide - LocalConnect Worker Features

## ✅ What's Been Implemented

### **Backend (New Endpoints)**
- ✅ `GET /profiles/me` - Get current worker's profile
- ✅ `GET /profiles/{id}/media` - List all work images/videos
- ✅ `PUT /profiles/{id}` - Update existing profile

### **Frontend (New Features)**
- ✅ Professional navbar with profile icon
- ✅ Separate login/logout buttons
- ✅ Worker-specific dashboard
- ✅ Automatic profile loading
- ✅ View/Edit profile modes
- ✅ Work image gallery
- ✅ Notification alerts for client requests
- ✅ Statistics overview

---

## 🎬 How to Test

### **Step 1: Open Application**
Open browser: `http://localhost:5173`

### **Step 2: Test as CLIENT**

1. **Login** with client account:
   ```
   Email: sarah@example.com
   Password: password
   ```

2. **After login**, you should:
   - Stay on **Home page** (`/`)
   - See **worker profiles** listed
   - See your name and avatar in navbar
   - See **Logout** button (red)

3. **Click on a worker** to view their profile

4. **Send contact request** to test notifications

5. **Logout** using the red button

---

### **Step 3: Test as WORKER**

1. **Logout** from client account (if logged in)

2. **Login** with worker account:
   ```
   Email: robert@example.com
   Password: password
   ```
   *(Note: If this doesn't work, the password might have been changed. Try registering a new worker account)*

3. **After login**, you should:
   - Be redirected to **Dashboard** (`/dashboard`)
   - See **4 stat cards** (Views, Requests, Rating, Earnings)
   - See your **existing profile** in view mode
   - See **"Edit Profile"** button
   - See **work gallery** section
   - See **contact requests** at bottom

4. **If you see notification banner**:
   - Orange/red alert at top
   - Shows number of client requests
   - Click "View Requests" to scroll down

5. **Test Profile Editing**:
   - Click **"Edit Profile"** button
   - Modify any field (e.g., change rate to $200)
   - Click **"Update Profile"**
   - Profile updates and switches back to view mode

6. **Test Image Upload**:
   - Scroll to **"Your Work Gallery"**
   - Click **"Choose File"**
   - Select an image from your computer
   - Click **"Upload Media"**
   - Image should appear in gallery below

7. **Check Navbar**:
   - See your **name and avatar** (👷 emoji)
   - See **"Worker"** role label
   - See **Dashboard** link in nav
   - See **Messages** link in nav
   - See **Logout** button (red)
   - Click on your **profile/avatar** to see dropdown

---

### **Step 4: Test New Worker Registration**

1. **Logout** from worker account

2. **Register** a new worker:
   ```
   Name: Test Worker
   Email: newworker@test.com
   Phone: +1-555-1234
   Password: test123
   Role: Worker (select from dropdown)
   ```

3. **After registration**, you should:
   - Be redirected to `/dashboard`
   - See **"Welcome! Let's create your profile"** message
   - See **profile creation form**
   - Fill out all fields
   - Click **"Create Profile"**
   - Profile created and switches to view mode

4. **Upload first work image**:
   - Gallery section becomes visible
   - Upload an image
   - See it appear in gallery

---

## 🔍 What to Verify

### **Navbar Behavior**
- [ ] **Not logged in**: Shows "Login" button (green)
- [ ] **Logged in as client**: Shows profile icon with name, "Logout" button
- [ ] **Logged in as worker**: Shows profile icon with 👷 emoji, "Logout" button, Dashboard link visible
- [ ] **Profile dropdown**: Click name/avatar to see dropdown with user info and quick links
- [ ] **Logout**: Clicking logout clears token and redirects to `/auth`

### **Worker Dashboard**
- [ ] **Profile auto-loads** when worker logs in
- [ ] **Stats display** correctly (views, requests, rating, earnings)
- [ ] **Notification banner** appears if contact requests exist
- [ ] **Profile in view mode** initially (if exists)
- [ ] **Edit button** switches to edit mode
- [ ] **Cancel button** switches back to view mode
- [ ] **Update saves** and refreshes profile
- [ ] **Work gallery** displays uploaded images in grid
- [ ] **Upload button** uploads and refreshes gallery
- [ ] **Contact requests** show with proper formatting

### **Client Home Page**
- [ ] **Worker profiles** load and display
- [ ] **Filter by category** works
- [ ] **Filter by location** works
- [ ] **View profile** navigates to detail page
- [ ] **Send contact** creates notification for worker

---

## 🐛 Troubleshooting

### **Issue: Worker sees "Welcome! Let's create your profile" but they already have one**
- **Fix**: Check browser console for API errors
- **Check**: `GET /profiles/me` endpoint response
- **Verify**: Worker is logged in with correct token

### **Issue: Images don't appear in gallery**
- **Fix**: Check `http://localhost:8080` is accessible
- **Verify**: Images uploaded successfully (check success message)
- **Check**: `/media` folder in backend contains files

### **Issue: No notifications appear**
- **Fix**: Login as **client** first, send contact request to worker
- **Then**: Login as **worker** to see notification
- **Verify**: Contact requests exist in database

### **Issue: Profile won't save**
- **Fix**: Check all required fields are filled
- **Verify**: Category is selected (not empty)
- **Check**: Backend logs for errors

### **Issue: Logout doesn't work**
- **Fix**: Hard refresh browser (`Cmd+Shift+R`)
- **Clear**: `localStorage` in browser DevTools
- **Verify**: `lc_token` is cleared after logout

---

## 🎯 Expected User Experience

### **New Worker Flow**
1. Register → Dashboard
2. See "Create Profile" prompt
3. Fill form → Save
4. Upload work images
5. Wait for client requests

### **Existing Worker Flow**
1. Login → Dashboard
2. See notification banner (if requests)
3. View stats and profile
4. Edit profile or upload images
5. Respond to client requests

### **Client Flow**
1. Login → Home page
2. Browse workers
3. Filter by category/location
4. View worker details
5. Send contact request

---

## 📝 Demo Credentials

### **Workers** (go to Dashboard)
- `robert@example.com` - Carpenter
- `jennifer@example.com` - Painter
- `maria@example.com` - Cook

### **Clients** (go to Home)
- `sarah@example.com` - Client
- `michael@example.com` - Client

**All passwords**: `password`

*(If password doesn't work, register a new account)*

---

## 🎨 Visual Checklist

When testing, verify these visual elements:

### **Navbar**
- [ ] Glassmorphism effect
- [ ] Profile icon circular with gradient
- [ ] Logout button red with hover effect
- [ ] Dropdown menu slides down smoothly

### **Dashboard**
- [ ] Stats cards colorful with gradients
- [ ] Notification banner orange/red with bell icon
- [ ] Profile view mode clean and organized
- [ ] Edit mode shows all form fields
- [ ] Gallery grid responsive
- [ ] Images have hover zoom effect
- [ ] Contact requests have avatar icons

### **Animations**
- [ ] Dropdown slides down
- [ ] Cards have hover effects
- [ ] Gallery images zoom on hover
- [ ] Buttons have transition effects

---

Ready to test! 🚀
