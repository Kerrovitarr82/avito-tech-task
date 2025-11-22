import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
    stages: [
        { duration: '30s', target: 5 },   // Разогрев до 5 RPS
        { duration: '1m', target: 5 },    // Держим 5 RPS
        { duration: '30s', target: 10 },  // Увеличиваем до 10 RPS
        { duration: '1m', target: 10 },   // Держим 10 RPS
        { duration: '30s', target: 0 },   // Остывание
    ],
    thresholds: {
        http_req_duration: ['p(95)<300'], // 95% запросов быстрее 300ms
        http_req_failed: ['rate<0.001'],  // Менее 0.1% ошибок (99.9% успешности)
        errors: ['rate<0.001'],
    },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    const scenarios = [
        // Сценарий 1: Получение статистики пользователя
        () => {
            const userId = `u${Math.floor(Math.random() * 200) + 1}`;
            const res = http.get(`${BASE_URL}/users/getStats?user_id=${userId}`);

            check(res, {
                'status is 200': (r) => r.status === 200,
                'response time < 300ms': (r) => r.timings.duration < 300,
            }) || errorRate.add(1);
        },

        // Сценарий 2: Получение PR для ревью
        () => {
            const userId = `u${Math.floor(Math.random() * 200) + 1}`;
            const res = http.get(`${BASE_URL}/users/getReview?user_id=${userId}`);

            check(res, {
                'status is 200': (r) => r.status === 200,
                'response time < 300ms': (r) => r.timings.duration < 300,
            }) || errorRate.add(1);
        },

        // Сценарий 3: Получение команды
        () => {
            const teamNum = Math.floor(Math.random() * 20) + 1;
            const res = http.get(`${BASE_URL}/team/get?team_name=Team-${teamNum}`);

            check(res, {
                'status is 200': (r) => r.status === 200,
                'response time < 300ms': (r) => r.timings.duration < 300,
            }) || errorRate.add(1);
        },

        // Сценарий 4: Создание PR
        () => {
            const prId = `pr-test-${Date.now()}-${Math.random()}`;
            const userId = `u${Math.floor(Math.random() * 200) + 1}`;

            const payload = JSON.stringify({
                pull_request_id: prId,
                pull_request_name: 'Load Test PR',
                author_id: userId,
            });

            const res = http.post(`${BASE_URL}/pullRequest/create`, payload, {
                headers: { 'Content-Type': 'application/json' },
            });

            check(res, {
                'status is 201': (r) => r.status === 201,
                'response time < 300ms': (r) => r.timings.duration < 300,
            }) || errorRate.add(1);
        },

        // Сценарий 5: Merge PR
        () => {
            const prNum = Math.floor(Math.random() * 50) + 1;
            const payload = JSON.stringify({
                pull_request_id: `pr-${prNum}`,
            });

            const res = http.post(`${BASE_URL}/pullRequest/merge`, payload, {
                headers: { 'Content-Type': 'application/json' },
            });

            check(res, {
                'status is 200': (r) => r.status === 200,
                'response time < 300ms': (r) => r.timings.duration < 300,
            }) || errorRate.add(1);
        },
    ];

    // Выбираем случайный сценарий
    const scenario = scenarios[Math.floor(Math.random() * scenarios.length)];
    scenario();

    sleep(0.1); // Небольшая пауза между запросами
}