% Read and analyze sweep data
% Initialize arrays to store results
results = struct('omega', [], 'th1max', [], 'th2max', []);

% Get list of CSV files in sweep directory
files = dir('./sweep/*.csv');
k = 1;

% Process each file
for i = 1:length(files)
    % Read CSV file
    filename = fullfile(files(i).folder, files(i).name);
    opts = detectImportOptions(filename);
    opts.DataLines = [1 inf];
    opts.VariableNames = {'time', 'th1', 'th2', 'th1dot', 'th2dot'};
    data = readtable(filename, opts);
    
    % Extract omega from filename
    [~, name, ~] = fileparts(files(i).name);
    parts = split(name, '_');
    omega = str2double(parts{end});
    
    % Find max in last 25% of data
    N = height(data);
    start_idx = floor(N * 0.75);
    th1max = max(data.th1(start_idx:end));
    th2max = max(data.th2(start_idx:end));
    
    % Store results
    results.omega(k) = omega;
    results.th1max(k) = th1max;
    results.th2max(k) = th2max;
    
    % Plot time series for k=20
    if k == 20
        figure;
        plot(data.time, data.th1, 'DisplayName', sprintf('omega=%.3f', omega));
        xlabel('Time');
        ylabel('theta1');
        title('Time Series for k=20');
        legend('show');
        grid on;
    end
    
    k = k + 1;
end

% Sort results by omega
[~, sort_idx] = sort(results.omega);
results.omega = results.omega(sort_idx);
results.th1max = results.th1max(sort_idx);
results.th2max = results.th2max(sort_idx);

% Plot frequency response
figure;
semilogx(results.omega, results.th1max, '.-', 'DisplayName', 'theta1');
hold on;
semilogx(results.omega, results.th2max, '.-', 'DisplayName', 'theta2');
xline(sqrt(0.382), 'k--', 'LineWidth', 0.5);
xline(sqrt(2.618), 'k--', 'LineWidth', 0.5);
xlabel('\omega');
ylabel('Amplitude');
title('Frequency Response');
legend('show');
grid on; 