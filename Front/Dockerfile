# Stage 1: Angular
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

#Сборка
RUN npm run build -- --configuration=production

# Stage 2: Nginx
FROM nginx:alpine

# Копия Nginx конфига
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копия из билда 
COPY --from=build /app/dist/FreeLib/browser /usr/share/nginx/html

# Открываем 80 порт
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
