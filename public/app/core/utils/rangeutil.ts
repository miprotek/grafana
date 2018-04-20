import _ from 'lodash';
import moment from 'moment';
import * as dateMath from './datemath';

var spans = {
  s: { display: 'Sekunde' },
  m: { display: 'Minute' },
  h: { display: 'Stunde' },
  d: { display: 'Tag' },
  w: { display: 'Woche' },
  M: { display: 'Monat' },
  y: { display: 'Jahr' },
};

var rangeOptions = [
  { from: 'now/d', to: 'now/d', display: 'Heute', section: 2 },
  { from: 'now/d', to: 'now', display: 'Heute bis jetzt', section: 2 },
  { from: 'now/w', to: 'now/w', display: 'Diese Woche', section: 2 },
  { from: 'now/w', to: 'now', display: 'Diese Woche bis jetzt', section: 2 },
  { from: 'now/M', to: 'now/M', display: 'Diesen Monat', section: 2 },
  { from: 'now/M', to: 'now', display: 'Diesen Monat bis jetzt', section: 2 },
  { from: 'now/y', to: 'now/y', display: 'Dieses Jahr', section: 2 },
  { from: 'now/y', to: 'now', display: 'Dieses Jahr bis jetzt', section: 2 },

  { from: 'now-1d/d', to: 'now-1d/d', display: 'Gestern', section: 1 },
  {
    from: 'now-2d/d',
    to: 'now-2d/d',
    display: 'Vorgestern',
    section: 1,
  },
  {
    from: 'now-7d/d',
    to: 'now-7d/d',
    display: 'Heute vor einer Woche',
    section: 1,
  },
  { from: 'now-1w/w', to: 'now-1w/w', display: 'Letzte Woche', section: 1 },
  { from: 'now-1M/M', to: 'now-1M/M', display: 'Letzten Monat', section: 1 },
  { from: 'now-1y/y', to: 'now-1y/y', display: 'Letztes Jahr', section: 1 },

  { from: 'now-5m', to: 'now', display: 'Letzten 5 Minuten', section: 3 },
  { from: 'now-15m', to: 'now', display: 'Letzten 15 Minuten', section: 3 },
  { from: 'now-30m', to: 'now', display: 'Letzten 30 Minuten', section: 3 },
  { from: 'now-1h', to: 'now', display: 'Letzte 1 Stunde', section: 3 },
  { from: 'now-3h', to: 'now', display: 'Letzten 3 Stunden', section: 3 },
  { from: 'now-6h', to: 'now', display: 'Letzten 6 Stunden', section: 3 },
  { from: 'now-12h', to: 'now', display: 'Letzten 12 Stunden', section: 3 },
  { from: 'now-24h', to: 'now', display: 'Letzten 24 Stunden', section: 3 },

  { from: 'now-2d', to: 'now', display: 'Letzten 2 Tage', section: 0 },
  { from: 'now-7d', to: 'now', display: 'Letzten 7 Tage', section: 0 },
  { from: 'now-30d', to: 'now', display: 'Letzten 30 Tage', section: 0 },
  { from: 'now-90d', to: 'now', display: 'Letzten 90 Tage', section: 0 },
  { from: 'now-6M', to: 'now', display: 'Letzten 6 Monate', section: 0 },
  { from: 'now-1y', to: 'now', display: 'Letztes 1 Jahr', section: 0 },
  { from: 'now-2y', to: 'now', display: 'Letzten 2 Jahre', section: 0 },
  { from: 'now-5y', to: 'now', display: 'Letzten 5 Jahre', section: 0 },
];

var absoluteFormat = 'D.MMM.YYYY HH:mm:ss';

var rangeIndex = {};
_.each(rangeOptions, function(frame) {
  rangeIndex[frame.from + ' bis ' + frame.to] = frame;
});

export function getRelativeTimesList(timepickerSettings, currentDisplay) {
  var groups = _.groupBy(rangeOptions, (option: any) => {
    option.active = option.display === currentDisplay;
    return option.section;
  });

  // _.each(timepickerSettings.time_options, (duration: string) => {
  //   let info = describeTextRange(duration);
  //   if (info.section) {
  //     groups[info.section].push(info);
  //   }
  // });

  return groups;
}

function formatDate(date) {
  return date.format(absoluteFormat);
}

// handles expressions like
// 5m
// 5m to now/d
// now/d to now
// now/d
// if no to <expr> then to now is assumed
export function describeTextRange(expr: any) {
  let isLast = expr.indexOf('+') !== 0;
  if (expr.indexOf('now') === -1) {
    expr = (isLast ? 'now-' : 'now') + expr;
  }

  let opt = rangeIndex[expr + ' bis jetzt'];
  if (opt) {
    return opt;
  }

  if (isLast) {
    opt = { from: expr, to: 'now' };
  } else {
    opt = { from: 'now', to: expr };
  }

  let parts = /^now([-+])(\d+)(\w)/.exec(expr);
  if (parts) {
    let unit = parts[3];
    let amount = parseInt(parts[2]);
    let span = spans[unit];
    if (span) {
      opt.display = isLast ? 'Last ' : 'Next ';
      opt.display += amount + ' ' + span.display;
      opt.section = span.section;
      if (amount > 1) {
        opt.display += 's';
      }
    }
  } else {
    opt.display = opt.from + ' bis ' + opt.to;
    opt.invalid = true;
  }

  return opt;
}

export function describeTimeRange(range) {
  var option = rangeIndex[range.from.toString() + ' bis ' + range.to.toString()];
  if (option) {
    return option.display;
  }

  if (moment.isMoment(range.from) && moment.isMoment(range.to)) {
    return formatDate(range.from) + ' bis ' + formatDate(range.to);
  }

  if (moment.isMoment(range.from)) {
    var toMoment = dateMath.parse(range.to, true);
    return formatDate(range.from) + ' bis ' + toMoment.fromNow();
  }

  if (moment.isMoment(range.to)) {
    var from = dateMath.parse(range.from, false);
    return from.fromNow() + ' bis ' + formatDate(range.to);
  }

  if (range.to.toString() === 'now') {
    var res = describeTextRange(range.from);
    return res.display;
  }

  return range.from.toString() + ' bis ' + range.to.toString();
}
