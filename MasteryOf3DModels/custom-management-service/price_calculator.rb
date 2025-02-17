class PriceCalculator
  def self.calculate(price, custom_params)
    cost_final = price # The initial price is the model's price

    custom_params.each do |key, value|
      # Adjustment for filament quality
      if key == 'filament_quality'
        case value
        when 'Low'
          cost_final += 0
        when 'Medium'
          cost_final += 5.0
        when 'High'
          cost_final += 10.0
        end
      end

      # Adjustment for color
      if key == 'color'
        case value
        when 'Red'
          cost_final += 0
        when 'Gold'
          cost_final += 5.0
        when 'Blue'
          cost_final += 3.0
        end
      end

      # Adjustment for size
      if key == 'size'
        case value
        when 'Small'
          cost_final += 0
        when 'Medium'
          cost_final += 5.0
        when 'Large'
          cost_final += 8.0
        end
      end

      # Adjustment for shape complexity
      if key == 'shape_complexity'
        case value
        when 'Low'
          cost_final += 0
        when 'Medium'
          cost_final += 6.0
        when 'High'
          cost_final += 12.0
        end
      end

      # Adjustment for material type
      if key == 'material_type'
        case value
        when 'Standard'
          cost_final += 0
        when 'Premium'
          cost_final += 15.0
        when 'Luxury'
          cost_final += 25.0
        end
      end
    end

    return cost_final
  end
end