# Go To-Do App

A simple To-Do application built using Go, Fiber, and MongoDB. This app allows you to create, read, update, and delete to-do items. The application exposes RESTful endpoints for interacting with the to-do list.

## Features
- **Create a new to-do**
- **View all to-dos**
- **Update a to-do** (requires `:id` in the URL)
- **Delete a to-do** (requires `:id` in the URL)

## Deployed Endpoint

The app is deployed on Railway. You can access the live API here:

**[https://gotodo-production-98ae.up.railway.app/todos](https://gotodo-production-98ae.up.railway.app/todos)**

## Endpoints

### 1. **GET /todos**
Fetch all the to-do items.

**Response Example:**
```json
[
    {
        "_id": "60f72b8f9b3d1a3c4b8e4725",
        "completed": false,
        "body": "Learn Go"
    },
    {
        "_id": "60f72b9f9b3d1a3c4b8e4726",
        "completed": true,
        "body": "Build an app"
    }
]
```
### 2. **POST /todos**
Create a new to-do item.

**Request Body:**
```json
{
    "body": "Learn Fiber",
    "completed": false
}
```
### 3. **PATCH /todos/:id**
Update a specific to-do item.

**Request Body:**
```json
{
    "completed": true
}
```
### 4. **Delete /todos/:id**
Delete a specific to-do item.