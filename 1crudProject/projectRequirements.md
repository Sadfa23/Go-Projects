# Project Requirements: Tasks API

## Core Features

1. **Create Task**  
    - Fields:  
      - `title` (required)  
      - `description`  
      - `due_date`  
      - `priority` (low/medium/high)  
      - `tags` (array)  
      - `completed` (boolean)  
      - `created_at`  
      - `updated_at`  

2. **Read Tasks**  
    - Retrieve:  
      - Single task  
      - List of tasks  
    - Filters:  
      - By `completed` status  
      - By `tag`  
      - By `due_date` range  
    - Sorting:  
      - By `due_date`  
      - By `priority`  

3. **Update Task**  
    - Update any field.  
    - Mark as complete/incomplete.  

4. **Delete Task**  
    - Choose one:  
      - Soft-delete (archived).  
      - Hard delete.  

5. **Validation**  
    - Basic validation: `title` is required.  

---

## Optional Extras

- Pagination for task listing.  
- Full-text search on `title` and `description`.  
- User authentication (JWT) to associate tasks with users.  
- Dockerfile and `docker-compose` setup with MongoDB.  
- Unit and integration tests.  

---

## Implementation Notes

- Decide early between soft-delete and hard delete.  
- For the initial build:  
  - Implement **hard delete**.  
  - Add **soft-delete** functionality later.  
Requirements & features (Tasks API)
