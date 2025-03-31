import { Pipe, PipeTransform } from '@angular/core';
import { Condition } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';

@Pipe({
  name: 'condition',
})
export class ActionConditionPipe implements PipeTransform {
  transform(condition: Condition | undefined): string {
    return parseCondition(condition);
  }
}
export function parseCondition(condition: Condition | undefined): string {
  const conditionType: Condition['conditionType']['value'] = condition?.conditionType.value;

  if (!conditionType) {
    return '';
  }

  if ('condition' in conditionType) {
    const { condition: innerCondition } = conditionType;

    if (typeof innerCondition.value === 'string') {
      // Applies for service, method condition of Request/ResponseCondition, event, and group of EventCondition
      return `${innerCondition.case}: ${innerCondition.value}`;
    }

    if ('all' in innerCondition) {
      // Applies for "all" condition of Request/ResponseCondition and EventCondition
      return condition?.conditionType.case ? `${condition?.conditionType.case}: all` : 'all';
    }
  }

  if ('name' in conditionType) {
    // Applies for function condition
    return `function: ${conditionType.name}`;
  }

  return '';
}
