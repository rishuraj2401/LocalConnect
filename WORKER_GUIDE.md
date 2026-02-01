# 👷 Worker Dashboard Guide

## Overview
Workers get a **dedicated dashboard** with profile management, work gallery, and client notifications.

---

## 🎯 What's Different for Workers vs Clients?

### **Workers** (Service Providers)
- **After Login** → Redirected to `/dashboard`
- **See**:
  - Their own profile (view/edit mode)
  - Statistics (views, requests, ratings, earnings)
  - **Notifications** for client contact requests
  - Work image/video gallery
  - Upload new work samples

### **Clients** (Service Seekers)
- **After Login** → Redirected to `/` (home page)
- **See**:
  - Browse all workers
  - Filter by category/location
  - View worker profiles
  - Send contact requests
  - Message workers

---

## 🚀 Worker Dashboard Features

### 1. **Profile Management**
- **View Mode**: See your complete profile with stats
- **Edit Mode**: Click "Edit Profile" button
- **First Time**: Automatically shown profile creation form

**Profile Fields**:
- Category (Carpenter, Painter, Cook, etc.)
- Location (City, State)
- Daily Rate ($)
- Years of Experience
- Bio (describe your skills)

### 2. **Notifications System** 🔔
When clients contact you, you'll see:
- **Big alert banner** at top of dashboard
- **Number of new requests** in red badge
- **"View Requests" button** to scroll to requests section

### 3. **Work Gallery** 📸
- **Upload** images or videos of your work
- **Showcase** your skills to attract clients
- **Gallery view** with hover effects
- Supports: JPG, PNG, MP4, etc.

### 4. **Contact Requests** 📬
See all client requests:
- Client's message
- Phone number (if shared)
- Date/time of request
- Clean card design with icons

### 5. **Statistics Dashboard** 📊
- **Profile Views** - How many people viewed your profile
- **Total Requests** - Number of clients who contacted you
- **Your Rating** - Average rating from reviews
- **Potential Earnings** - Estimated earnings from requests

---

## 🔐 Testing the Worker Experience

### **Step 1: Login as a Worker**
Use one of these demo accounts (password: `password`):
```
robert@example.com    - Carpenter, New York
jennifer@example.com  - Painter, Los Angeles
maria@example.com     - Cook, Houston
```

### **Step 2: View Your Dashboard**
After login, you'll be redirected to `/dashboard` where you'll see:
1. **Welcome message** (if no profile)
2. **Stats overview** (if profile exists)
3. **Notification banner** (if contact requests)
4. **Your profile** (view or edit mode)
5. **Work gallery** (empty or with images)
6. **Contact requests** list

### **Step 3: Test Features**

#### **Upload Work Images**
1. Scroll to "Your Work Gallery"
2. Click "Choose File" and select an image
3. Click "Upload Media"
4. Image appears in gallery below

#### **Edit Your Profile**
1. Click "Edit Profile" button
2. Modify any fields
3. Click "Update Profile"
4. Profile updates instantly

#### **View Notifications**
1. See notification banner at top
2. Click "View Requests" button
3. Scroll to Contact Requests section
4. See all client messages with details

---

## 🎨 UI Design Highlights

### **Navbar (when logged in as worker)**
```
[LocalConnect] | Find Workers | Messages | Dashboard | [👷 Your Name] [Logout]
                                                         Worker
```

- **Profile card** shows your avatar, name, and role
- **Logout button** is always visible (red)
- **Click profile** to see dropdown with quick links

### **Dashboard Sections**
1. **Hero** - Title and description
2. **Notifications** - Orange/red alert banner (if requests)
3. **Stats** - 4 colorful cards with key metrics
4. **Profile** - View/edit your professional info
5. **Gallery** - Grid of work images/videos
6. **Requests** - List of client contact requests

### **Color Coding**
- **Green/Teal** - Primary actions, profile
- **Orange/Red** - Notifications, alerts
- **Blue** - Secondary info
- **Gold** - Ratings

---

## 🔄 Login Flow

### **Worker Login**
```
/auth (login) → [Success] → /dashboard
                           ↓
                     Load profile data
                           ↓
                  Show dashboard or setup
```

### **Client Login**
```
/auth (login) → [Success] → / (home)
                           ↓
                     Browse workers
```

---

## 📱 Responsive Design
- Works on desktop, tablet, and mobile
- Glassmorphism effects
- Smooth animations
- Professional gradient backgrounds

---

## 🛠️ API Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/profiles/me` | GET | Get current worker's profile |
| `/profiles` | POST | Create new profile |
| `/profiles/{id}` | PUT | Update profile |
| `/profiles/{id}/media` | GET | List work images |
| `/profiles/{id}/media` | POST | Upload work image |
| `/contact-requests` | GET | Get client requests |
| `/categories` | GET | Get service categories |

---

## ✅ Ready to Use!

1. **Backend running**: `./start.sh` in `backend/` folder
2. **Frontend running**: `npm run dev` in `frontend/` folder
3. **Demo data loaded**: `./load-demo-data.sh`
4. **Login as worker**: Use demo credentials above
5. **Enjoy the new dashboard!** 🎉

---

## 🆕 New Features Summary

✅ **Separate worker page** - Dedicated dashboard for workers  
✅ **Profile display** - Workers see their own profile  
✅ **Notifications** - Alert banner for new client requests  
✅ **Image uploads** - Workers can upload work samples  
✅ **Work gallery** - Display uploaded images in grid  
✅ **Edit profile** - Toggle between view/edit modes  
✅ **Statistics** - Key metrics at a glance  
✅ **Clean navbar** - Profile icon with logout button  

---

Need help? Check the browser console or backend logs for any errors!
